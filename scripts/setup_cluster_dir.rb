#!/usr/bin/ruby

require 'optparse'
require 'fileutils'
require 'json'

# nodeid_prefix defines the prefix of each node in
# the cluster. Default value is "node"
nodeid_prefix = "node"

# rpcport_start specifies the start of the port range
# for RPC servers. Default value is 6666. If the cluster
# size is 3, then rpc ports will be 6666, 6667 and 6668
rpcport_start = 6666

# apiport_start specifies the start of the port range
# for API servers. Default value is 7777. If the cluster
# size is 3, then api ports will be 7777, 7778, 7779
apiport_start = 7777

# cluster_size specifies the number of nodes in the cluster
# It is recommended to have odd number of nodes
cluster_size = 3

# cluster_dir specifies the root of the directory where
# all cluster data is stored for local testing
cluster_dir = './cluster'

OptionParser.new do |opts|
    opts.banner = "Usage: setup_cluster_dir.rb [options]"

    opts.on('-p', '--nodeid-prefix', 'Prefix for generated node ID') { |v| nodeid_prefix = v }
    opts.on('-c', '--rpcport-start', 'Start of rpc-server port range') { |v| rpcport_start = v.to_i }
    opts.on('-a', '--apiport-start', 'Start of api-server port range') { |v| apiport_start = v.to_i }
    opts.on('-s', '--cluster-size', 'Size of the cluster') { |v| cluster_size = v.to_i }
    opts.on('-d', '--cluster-dir', 'Root directory for cluster') { |v| cluster_dir = v }
end.parse!

# NodeInfo represents the information on a particular node in the
# raft cluster like its name, url of RPC and API servers. This can
# also be serialized and deserialized to/from JSON
class NodeInfo
    attr_accessor :node_id, :rpc_url, :api_url
    
    def initialize(node_id, rpc_port, api_port)
        @node_id = node_id
        @rpc_url = "localhost:#{rpc_port}"
        @api_url = "localhost:#{api_port}"
    end
    
    # Convert NodeInfo to JSON format
    def to_json(*a)
        {
            node_id: @node_id, 
            rpc_url: @rpc_url, 
            api_url: @api_url
        }.to_json(*a)
    end

    # Parse JSON and try to get NodeInfo
    def self.from_json string
        data = JSON.load string
        self.new data['node_id'], data['rpc_url'], data['api_url'] 
    end
end

# Generate node information in JSON from given parameters
cluster_node_info = (1..cluster_size).map do |serial_no|
    cur_node_id = "#{nodeid_prefix}-#{serial_no}"
    cur_rpc_port = rpcport_start + (serial_no - 1)
    cur_api_port = apiport_start + (serial_no - 1)
    NodeInfo.new cur_node_id, cur_rpc_port, cur_api_port
end

json_cluster_node_info = cluster_node_info.to_json

# Generate directory structure for each node. The entire cluster
# information will be kept in cluster/ directory. It will have a
# directory for each node with node_id being its name. It will have
# subdirectory for cluster configuration, entry, metadata, state
# and snapshot persistence.
cluster_node_info.each do |node_info|
    cur_node_id = node_info.node_id
    puts "Generating directory for #{cur_node_id}"
    node_dir = "#{cluster_dir}/#{cur_node_id}"
    node_directories = {
        cluster_config_dir: "#{node_dir}/cluster",
        entry_dir:          "#{node_dir}/data/entry",
        metadata_dir:       "#{node_dir}/data/metadata",
        state_dir:          "#{node_dir}/state",
        snapshot_dir:       "#{node_dir}/snapshot"
    }
    # Generate directory structure for the current node
    node_directories.each do |dir_purpose, dir_path|
        puts "* Generating #{dir_path}"
        FileUtils.mkdir_p dir_path
    end

    # Write cluster configuration file to the appropriate directory
    # Name of the file should be config.json
    cluster_config_file = "#{node_directories[:cluster_config_dir]}/config.json"
    File.open(cluster_config_file, 'w+') do |f|
        f.write json_cluster_node_info
    end
end