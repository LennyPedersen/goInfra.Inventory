package system

import (
	"errors"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func GetLocalIP() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "", err
    }

    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String(), nil
            }
        }
    }
    return "", errors.New("no IP address found")
}

func GetHostname() (string, error) {
    hostname, err := exec.Command("hostname").Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(hostname)), nil
}

func GetServices() (string, error) {
    var services []string
    switch runtime.GOOS {
    case "windows":
        services = checkWindowsServices()
    case "linux":
        services = checkLinuxServices()
    case "darwin":
        services = checkMacServices()
    default:
        return "", errors.New("unsupported operating system")
    }
    return strings.Join(services, ", "), nil
}

func checkWindowsServices() []string {
    var services []string
    out, _ := exec.Command("powershell", "Get-Service").Output()
    if strings.Contains(string(out), "MSSQL") {
        services = append(services, "SQL Server")
    }
    if strings.Contains(string(out), "W3SVC") {
        services = append(services, "IIS")
    }
    if strings.Contains(string(out), "Docker") {
        services = append(services, "Docker")
    }
    return services
}

func checkLinuxServices() []string {
    var services []string
    out, _ := exec.Command("systemctl", "list-units", "--type=service").Output()
    if strings.Contains(string(out), "nginx") {
        services = append(services, "Nginx")
    }
    if strings.Contains(string(out), "docker") {
        services = append(services, "Docker")
    }
    if strings.Contains(string(out), "mysql") {
        services = append(services, "MySQL")
    }
    if strings.Contains(string(out), "postgresql") {
        services = append(services, "PostgreSQL")
    }
    return services
}

func checkMacServices() []string {
    var services []string
    out, _ := exec.Command("brew", "services", "list").Output()
    if strings.Contains(string(out), "nginx") {
        services = append(services, "Nginx")
    }
    if strings.Contains(string(out), "docker") {
        services = append(services, "Docker")
    }
    if strings.Contains(string(out), "mysql") {
        services = append(services, "MySQL")
    }
    if strings.Contains(string(out), "postgresql") {
        services = append(services, "PostgreSQL")
    }
    return services
}

func GetOSVersion() (string, error) {
    switch runtime.GOOS {
    case "windows":
        return getWindowsVersion()
    case "linux":
        return getLinuxVersion()
    case "darwin":
        return getMacVersion()
    default:
        return "", errors.New("unsupported operating system")
    }
}

func getWindowsVersion() (string, error) {
    out, err := exec.Command("cmd", "/c", "ver").Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}

func getLinuxVersion() (string, error) {
    out, err := exec.Command("lsb_release", "-a").Output()
    if err != nil {
        return "", err
    }
    return string(out), nil
}

func getMacVersion() (string, error) {
    out, err := exec.Command("sw_vers").Output()
    if err != nil {
        return "", err
    }
    return string(out), nil
}


func GetOpenPorts() (string, error) {
    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("powershell", "Get-NetTCPConnection", "|", "Select-Object", "LocalPort")
    case "linux":
        cmd = exec.Command("ss", "-tuln")
    case "darwin":
        cmd = exec.Command("netstat", "-an")
    default:
        return "", errors.New("unsupported operating system")
    }

    out, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}
