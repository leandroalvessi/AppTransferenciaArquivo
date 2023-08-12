package conf_rede

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

type NameIP struct {
	Name string
	IP   string
}

func GetLocalIP() ([]NameIP, error) {
	var nameIPs []NameIP

	// Obtém todas as interfaces de rede disponíveis
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Itera sobre as interfaces e encontra o endereço IP local
	for _, iface := range interfaces {
		// Ignora interfaces não ativas e loopback (127.0.0.1)
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Obtém os endereços IP da interface
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		// Itera sobre os endereços IP da interface e verifica se é um endereço IPv4
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				nameIP := NameIP{
					Name: iface.Name,
					IP:   ipnet.IP.String(),
				}
				nameIPs = append(nameIPs, nameIP)
			}
		}
	}

	if len(nameIPs) == 0 {
		return nil, fmt.Errorf("Endereço IP não encontrado")
	}

	return nameIPs, nil
}

func OpenBrowser(url string) error {
	var err error

	switch os := runtime.GOOS; os {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		err = fmt.Errorf("Plataforma não suportada: %s", os)
	}

	return err
}
