package utils

import ("os")

PATH := []string {"/c/Users/abhin/.nvm/versions/node/v22.12.0/bin", "/c/Users/abhin/bin", "/mingw64/bin", 
"/usr/local/bin", "/usr/bin", "/bin"}

func findFileInDir(target string) (string, bool, error) {
	for _, i := range PATH{
		entries, err := os.ReadDir(i)
		if err != nil{
			return "", false, err
		}

		for _, entry := range entries{
			if !entry.IsDir() && entry.Name() == target{
				return i + "/" + entry.Name(), true, nil
			}
		}
	}
	return "", false, nil
}

func main() {
	filePath, found, err := findFileInDir(".", "main.go")
	
}