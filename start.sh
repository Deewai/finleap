GREEN='\033[0;32m'
NC='\033[0m'
printf "\n${GREEN}Starting Application...\n"
printf "${NC}"
docker-compose up --build
printf "\n${GREEN}minikube ip is $(minikube ip)\n"
printf "\nApplication can be tested with 'curl $(minikube ip):80/tree -H Host:local.ecosia.org'\n"
printf "${NC}"