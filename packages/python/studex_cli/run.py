import os
import sys
import platform
import urllib.request
import tarfile
import zipfile
import subprocess

VERSION = "1.0.0"

def get_platform_info():
    sys_plat = platform.system().lower()
    machine = platform.machine().lower()
    
    if sys_plat == "windows":
        goos = "windows"
    elif sys_plat == "darwin":
        goos = "darwin"
    else:
        goos = "linux"
        
    if machine in ["x86_64", "amd64"]:
        goarch = "x86_64"
    elif machine in ["aarch64", "arm64"]:
        goarch = "arm64"
    elif machine in ["i386", "i686", "x86"]:
        goarch = "i386"
    else:
        goarch = machine
        
    return goos, goarch

def ensure_binary():
    goos, goarch = get_platform_info()
    ext = "zip" if goos == "windows" else "tar.gz"
    filename = f"studex-cli_{goos}_{goarch}.{ext}"
    
    base_dir = os.path.dirname(os.path.abspath(__file__))
    bin_dir = os.path.join(base_dir, "bin")
    os.makedirs(bin_dir, exist_ok=True)
    
    exe_ext = ".exe" if goos == "windows" else ""
    bin_path = os.path.join(bin_dir, f"studex-cli{exe_ext}")
    
    if os.path.exists(bin_path):
        return bin_path
        
    url = f"https://github.com/itshivams/Studex-CLI/releases/download/v{VERSION}/{filename}"
    archive_path = os.path.join(bin_dir, filename)
    
    print(f"Downloading studex-cli from {url}...")
    try:
        urllib.request.urlretrieve(url, archive_path)
    except Exception as e:
        print(f"Failed to download from {url}: {e}")
        sys.exit(1)
        
    print("Extracting...")
    try:
        if ext == "zip":
            with zipfile.ZipFile(archive_path, 'r') as zip_ref:
                zip_ref.extractall(bin_dir)
        else:
            with tarfile.open(archive_path, "r:gz") as tar_ref:
                tar_ref.extractall(path=bin_dir)
                
        if goos != "windows":
            os.chmod(bin_path, 0o755)
            
        os.remove(archive_path)
    except Exception as e:
        print(f"Failed to extract {archive_path}: {e}")
        sys.exit(1)
        
    return bin_path

def main():
    bin_path = ensure_binary()
    args = sys.argv[1:]
    
    try:
        result = subprocess.run([bin_path] + args)
        sys.exit(result.returncode)
    except Exception as e:
        print(f"Error executing studex-cli: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
