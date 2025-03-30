#!/usr/bin/env python3

import os
import subprocess
from pathlib import Path

def ensure_directory(path):
    """Create directory if it doesn't exist."""
    path.mkdir(parents=True, exist_ok=True)

def generate_proto(proto_file, output_dir):
    """Generate protobuf templates using protoc."""
    try:
        # Create output directory if it doesn't exist
        ensure_directory(output_dir)
        
        # Get the parent directory of the proto file for imports
        proto_dir = proto_file.parent.parent
        
        # Run protoc command
        cmd = [
            'protoc',
            f'--proto_path={proto_dir}',  # Add proto path for imports
            f'--go_out={output_dir}',
            f'--go-grpc_out={output_dir}',
            str(proto_file)
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode != 0:
            print(f"Error generating proto for {proto_file}:")
            print(result.stderr)
        else:
            print(f"Successfully generated proto for {proto_file}")
            
    except Exception as e:
        print(f"Error processing {proto_file}: {str(e)}")

def main():
    # Get current directory
    current_dir = Path.cwd()
    
    # Scan for folders containing api/<service>.proto
    for item in current_dir.iterdir():
        if not item.is_dir():
            continue
            
        # Check for api directory
        api_dir = item / 'api'
        if not api_dir.exists() or not api_dir.is_dir():
            continue
            
        # Look for .proto files in api directory
        for proto_file in api_dir.glob('*.proto'):
            # Get service name (filename without extension)
            service_name = proto_file.stem
            
            # Create output directory path
            output_dir = current_dir # / 'library' / 'proto' / service_name
            
            print(f"Processing {proto_file}")
            generate_proto(proto_file, output_dir)

if __name__ == '__main__':
    main()
