#!/bin/bash

# Script de démarrage pour FactuGest-WebInformatique
# Usage: ./start.sh [start|stop|restart|status]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
PID_FILE="$BACKEND_DIR/factugest.pid"
LOG_FILE="$BACKEND_DIR/factugest.log"

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction d'affichage avec couleurs
print_status() {
    echo -e "${BLUE}[FactuGest]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[FactuGest]${NC} $1"
}

print_error() {
    echo -e "${RED}[FactuGest]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[FactuGest]${NC} $1"
}

# Vérifier si le serveur est en cours d'exécution
is_running() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0
        else
            rm -f "$PID_FILE"
            return 1
        fi
    fi
    return 1
}

# Démarrer le serveur
start_server() {
    print_status "Démarrage de FactuGest-WebInformatique..."
    
    if is_running; then
        print_warning "Le serveur est déjà en cours d'exécution."
        return 1
    fi
    
    # Vérifier que MySQL est démarré
    if ! systemctl is-active --quiet mysql; then
        print_status "Démarrage de MySQL..."
        sudo systemctl start mysql
        sleep 2
    fi
    
    # Aller dans le répertoire backend
    cd "$BACKEND_DIR" || {
        print_error "Impossible d'accéder au répertoire backend"
        exit 1
    }
    
    # Vérifier que l'exécutable existe
    if [ ! -f "factugest-server" ]; then
        print_status "Compilation du serveur..."
        go build -o factugest-server . || {
            print_error "Erreur lors de la compilation"
            exit 1
        }
    fi
    
    # Démarrer le serveur en arrière-plan
    nohup ./factugest-server > "$LOG_FILE" 2>&1 &
    local pid=$!
    echo $pid > "$PID_FILE"
    
    # Attendre un peu et vérifier que le serveur a démarré
    sleep 3
    if is_running; then
        print_success "Serveur démarré avec succès (PID: $pid)"
        print_status "Interface web disponible sur : http://localhost:8080"
        print_status "Logs disponibles dans : $LOG_FILE"
    else
        print_error "Erreur lors du démarrage du serveur"
        if [ -f "$LOG_FILE" ]; then
            print_error "Dernières lignes du log :"
            tail -n 10 "$LOG_FILE"
        fi
        exit 1
    fi
}

# Arrêter le serveur
stop_server() {
    print_status "Arrêt de FactuGest-WebInformatique..."
    
    if ! is_running; then
        print_warning "Le serveur n'est pas en cours d'exécution."
        return 1
    fi
    
    local pid=$(cat "$PID_FILE")
    kill "$pid" 2>/dev/null
    
    # Attendre que le processus se termine
    local count=0
    while ps -p "$pid" > /dev/null 2>&1 && [ $count -lt 10 ]; do
        sleep 1
        count=$((count + 1))
    done
    
    if ps -p "$pid" > /dev/null 2>&1; then
        print_warning "Arrêt forcé du serveur..."
        kill -9 "$pid" 2>/dev/null
    fi
    
    rm -f "$PID_FILE"
    print_success "Serveur arrêté avec succès"
}

# Redémarrer le serveur
restart_server() {
    stop_server
    sleep 2
    start_server
}

# Afficher le statut
show_status() {
    print_status "Statut de FactuGest-WebInformatique :"
    
    if is_running; then
        local pid=$(cat "$PID_FILE")
        print_success "Serveur en cours d'exécution (PID: $pid)"
        print_status "Interface web : http://localhost:8080"
        
        # Vérifier la connectivité
        if command -v curl > /dev/null 2>&1; then
            if curl -s http://localhost:8080 > /dev/null; then
                print_success "Interface web accessible"
            else
                print_warning "Interface web non accessible"
            fi
        fi
    else
        print_error "Serveur arrêté"
    fi
    
    # Statut de MySQL
    if systemctl is-active --quiet mysql; then
        print_success "MySQL en cours d'exécution"
    else
        print_error "MySQL arrêté"
    fi
}

# Afficher les logs
show_logs() {
    if [ -f "$LOG_FILE" ]; then
        print_status "Dernières lignes du log :"
        tail -n 20 "$LOG_FILE"
    else
        print_warning "Aucun fichier de log trouvé"
    fi
}

# Fonction d'aide
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commandes disponibles :"
    echo "  start     Démarrer le serveur FactuGest"
    echo "  stop      Arrêter le serveur FactuGest"
    echo "  restart   Redémarrer le serveur FactuGest"
    echo "  status    Afficher le statut du serveur"
    echo "  logs      Afficher les derniers logs"
    echo "  help      Afficher cette aide"
    echo ""
    echo "Si aucune commande n'est spécifiée, 'start' est utilisé par défaut."
}

# Script principal
case "${1:-start}" in
    start)
        start_server
        ;;
    stop)
        stop_server
        ;;
    restart)
        restart_server
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Commande inconnue : $1"
        show_help
        exit 1
        ;;
esac

