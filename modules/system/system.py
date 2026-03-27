import os
import shutil
import time

def get_data(config):
    disks_input = config.get('disks', [("/", "OS")])
    results = []
    
    disks = []
    for item in disks_input:
        if isinstance(item, dict):
            disks.append((item.get('path', '/'), item.get('label', 'Disk')))
        elif isinstance(item, (list, tuple)) and len(item) >= 2:
            disks.append((item[0], item[1]))

    for path, label in disks:
        try:
            total, used, free = shutil.disk_usage(path)
            total_gb = total / (1024**3)
            used_gb = used / (1024**3)
            percent = (used / total) * 100
            results.append(f"Disk {label}: {used_gb:.1f} GB / {total_gb:.1f} GB ({percent:.1f}% used)")
        except Exception:
            results.append(f"Disk {path}: Not found")
            
    try:
        with open("/proc/meminfo", "r") as f:
            memlines = f.readlines()
            total_kb = int(memlines[0].split()[1])
            avail_kb = 0
            for line in memlines:
                if "MemAvailable" in line:
                    avail_kb = int(line.split()[1])
                    break
            if avail_kb == 0:
                avail_kb = int(memlines[1].split()[1])
                
            used_gb = (total_kb - avail_kb) / (1024 * 1024)
            total_gb = total_kb / (1024 * 1024)
            results.append(f"Ram usage: {used_gb:.1f} GB / {total_gb:.1f} GB")
    except Exception:
        results.append("Ram usage: Unavailable")
        
    try:
        def get_cpu_times():
            with open("/proc/stat", "r") as f:
                line = f.readline()
                fields = [float(column) for column in line.strip().split()[1:]]
                return sum(fields), fields[3]
        
        t1, i1 = get_cpu_times()
        time.sleep(0.1)
        t2, i2 = get_cpu_times()
        cpu_usage = 100 * (1 - (i2 - i1) / (t2 - t1))
        results.append(f"CPU Usage: {cpu_usage:.1f}%")
        
        temp = "N/A"
        for tzone in ['thermal_zone0', 'thermal_zone1']:
            tpath = f"/sys/class/thermal/{tzone}/temp"
            if os.path.exists(tpath):
                with open(tpath, "r") as f:
                    temp = f"{float(f.read()) / 1000.0:.1f}"
                    break
        results.append(f"CPU Temp: {temp}C")
    except Exception:
        results.append("CPU Stats: Unavailable")
        
    return results
