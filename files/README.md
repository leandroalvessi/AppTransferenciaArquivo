git config --global color.ui true
git config --global user.name "Leandro Alves da Silva"
git config --global user.email "leandroalves.desenvolvedor@gmail.com"
ssh-keygen -t ed25519 -C "leandroalves.desenvolvedor@gmail.com"

cat ~/.ssh/id_ed25519.pub

ssh -T git@github.com

Build
GOOS="windows" go build

Ver variais de ambiente
go env