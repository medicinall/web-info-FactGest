// Configuration de l'API
const API_BASE_URL = '/api';

// Variables globales
let currentUser = null;
let clients = [];
let factures = [];

// Utilitaires
function showLoading() {
    document.getElementById('loading-overlay').classList.add('show');
}

function hideLoading() {
    document.getElementById('loading-overlay').classList.remove('show');
}

function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.remove();
    }, 5000);
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('fr-FR');
}

function formatCurrency(amount) {
    return new Intl.NumberFormat('fr-FR', {
        style: 'currency',
        currency: 'EUR'
    }).format(amount);
}

// Fonction pour générer et télécharger un PDF de facture
async function generatePDF(factureId) {
    try {
        showLoading();
        showNotification('Génération du PDF en cours...', 'info');
        
        const response = await fetch(`${API_BASE_URL}/factures/${factureId}/pdf`);
        
        if (!response.ok) {
            throw new Error('Erreur lors de la génération du PDF');
        }
        
        // Récupérer le nom du fichier depuis les en-têtes
        const contentDisposition = response.headers.get('Content-Disposition');
        let filename = 'facture.pdf';
        if (contentDisposition) {
            const filenameMatch = contentDisposition.match(/filename="(.+)"/);
            if (filenameMatch) {
                filename = filenameMatch[1];
            }
        }
        
        // Créer un blob avec le contenu PDF
        const blob = await response.blob();
        
        // Créer un lien de téléchargement
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        
        // Nettoyer
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        
        showNotification('PDF téléchargé avec succès', 'success');
    } catch (error) {
        showNotification('Erreur lors de la génération du PDF: ' + error.message, 'error');
    } finally {
        hideLoading();
    }
}

// API Calls
async function apiCall(endpoint, options = {}) {
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Erreur API');
        }
        
        return data;
    } catch (error) {
        console.error('Erreur API:', error);
        throw error;
    }
}

// Gestion des onglets
function initTabs() {
    const tabButtons = document.querySelectorAll('.nav-tab');
    const tabContents = document.querySelectorAll('.tab-content');
    
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabName = button.dataset.tab;
            
            // Mettre à jour les boutons
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            // Mettre à jour le contenu
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`${tabName}-tab`).classList.add('active');
            
            // Charger les données si nécessaire
            loadTabData(tabName);
        });
    });
}

function loadTabData(tabName) {
    switch (tabName) {
        case 'factures':
            loadFactures();
            break;
        case 'clients':
            loadClients();
            break;
        case 'admin':
            loadAdminStats();
            break;
    }
}

// Gestion des clients
async function loadClients() {
    try {
        showLoading();
        const response = await apiCall('/clients');
        clients = response.data || [];
        renderClientsTable();
        updateClientSelect();
    } catch (error) {
        showNotification('Erreur lors du chargement des clients', 'error');
    } finally {
        hideLoading();
    }
}

function renderClientsTable() {
    const tbody = document.querySelector('#clients-table tbody');
    tbody.innerHTML = '';
    
    clients.forEach(client => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${client.nom}</td>
            <td>${client.prenom}</td>
            <td>${client.entreprise || '-'}</td>
            <td>${client.email || '-'}</td>
            <td>${client.telephone || '-'}</td>
            <td>${client.ville || '-'}</td>
            <td class="actions">
                <button class="action-btn edit" onclick="editClient(${client.id})" title="Modifier">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="action-btn delete" onclick="deleteClient(${client.id})" title="Supprimer">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

function updateClientSelect() {
    const select = document.getElementById('facture-client');
    select.innerHTML = '<option value="">Sélectionner un client</option>';
    
    clients.forEach(client => {
        const option = document.createElement('option');
        option.value = client.id;
        option.textContent = `${client.nom} ${client.prenom}${client.entreprise ? ' - ' + client.entreprise : ''}`;
        select.appendChild(option);
    });
}

async function saveClient(clientData) {
    try {
        showLoading();
        const response = await apiCall('/clients', {
            method: 'POST',
            body: JSON.stringify(clientData)
        });
        
        showNotification('Client créé avec succès', 'success');
        loadClients();
        resetClientForm();
    } catch (error) {
        showNotification(error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function updateClient(id, clientData) {
    try {
        showLoading();
        const response = await apiCall(`/clients/${id}`, {
            method: 'PUT',
            body: JSON.stringify(clientData)
        });
        
        showNotification('Client mis à jour avec succès', 'success');
        loadClients();
    } catch (error) {
        showNotification(error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function deleteClient(id) {
    if (!confirm('Êtes-vous sûr de vouloir supprimer ce client ?')) {
        return;
    }
    
    try {
        showLoading();
        await apiCall(`/clients/${id}`, { method: 'DELETE' });
        showNotification('Client supprimé avec succès', 'success');
        loadClients();
    } catch (error) {
        showNotification(error.message, 'error');
    } finally {
        hideLoading();
    }
}

function editClient(id) {
    const client = clients.find(c => c.id === id);
    if (!client) return;
    
    // Remplir le formulaire avec les données du client
    document.getElementById('client-nom').value = client.nom;
    document.getElementById('client-prenom').value = client.prenom;
    document.getElementById('client-entreprise').value = client.entreprise || '';
    document.getElementById('client-siret').value = client.siret || '';
    document.getElementById('client-adresse').value = client.adresse || '';
    document.getElementById('client-ville').value = client.ville || '';
    document.getElementById('client-code-postal').value = client.code_postal || '';
    document.getElementById('client-telephone').value = client.telephone || '';
    document.getElementById('client-email').value = client.email || '';
    
    // Changer vers l'onglet de création/modification
    document.querySelector('[data-tab="nouveau-client"]').click();
    
    // Modifier le formulaire pour la mise à jour
    const form = document.getElementById('client-form');
    form.dataset.editId = id;
    document.querySelector('#nouveau-client-tab .tab-header h2').innerHTML = 
        '<i class="fas fa-user-edit"></i> Modifier le client';
}

function resetClientForm() {
    const form = document.getElementById('client-form');
    form.reset();
    delete form.dataset.editId;
    document.querySelector('#nouveau-client-tab .tab-header h2').innerHTML = 
        '<i class="fas fa-user-plus"></i> Créer un nouveau client';
}

// Gestion des factures
async function loadFactures() {
    try {
        showLoading();
        const response = await apiCall('/factures');
        factures = response.data || [];
        renderFacturesTable();
    } catch (error) {
        showNotification('Erreur lors du chargement des factures', 'error');
    } finally {
        hideLoading();
    }
}

function renderFacturesTable() {
    const tbody = document.querySelector('#factures-table tbody');
    tbody.innerHTML = '';
    
    factures.forEach(facture => {
        const row = document.createElement('tr');
        const clientName = facture.client ? 
            `${facture.client.nom} ${facture.client.prenom}` : 
            'Client inconnu';
        
        row.innerHTML = `
            <td>${facture.numero_facture}</td>
            <td>${clientName}</td>
            <td>${formatDate(facture.date_facture)}</td>
            <td>${formatDate(facture.date_echeance)}</td>
            <td>${formatCurrency(facture.montant_ttc)}</td>
            <td><span class="status ${facture.statut}">${facture.statut}</span></td>
            <td class="actions">
                <button class="action-btn edit" onclick="editFacture(${facture.id})" title="Modifier">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="action-btn pdf" onclick="generatePDF(${facture.id})" title="PDF">
                    <i class="fas fa-file-pdf"></i>
                </button>
                <button class="action-btn delete" onclick="deleteFacture(${facture.id})" title="Supprimer">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

async function saveFacture(factureData) {
    try {
        showLoading();
        const response = await apiCall('/factures', {
            method: 'POST',
            body: JSON.stringify(factureData)
        });
        
        showNotification('Facture créée avec succès', 'success');
        loadFactures();
        closeModal('facture-modal');
    } catch (error) {
        showNotification(error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function updateFacture(id, factureData) {
    try {
        showLoading();
        const response = await apiCall(`/factures/${id}`, {
            method: 'PUT',
            body: JSON.stringify(factureData)
        });
        
        showNotification('Facture mise à jour avec succès', 'success');
        loadFactures();
        closeModal('facture-modal');
    } catch (error) {
        showNotification(error.message, 'error');
    } finally {
        hideLoading();
    }
}

async function deleteFacture(id) {
    if (!confirm('Êtes-vous sûr de vouloir supprimer cette facture ?')) {
        return;
    }
    
    try {
        showLoading();
        await apiCall(`/factures/${id}`, { method: 'DELETE' });
        showNotification('Facture supprimée avec succès', 'success');
        loadFactures();
    } catch (error) {
        showNotification(error.message, 'error');
    } finally {
        hideLoading();
    }
}

function editFacture(id) {
    const facture = factures.find(f => f.id === id);
    if (!facture) return;
    
    // Remplir le formulaire avec les données de la facture
    document.getElementById('facture-numero').value = facture.numero_facture;
    document.getElementById('facture-client').value = facture.client_id;
    document.getElementById('facture-date').value = facture.date_facture.split('T')[0];
    document.getElementById('facture-echeance').value = facture.date_echeance.split('T')[0];
    document.getElementById('facture-statut').value = facture.statut;
    document.getElementById('facture-montant-ht').value = facture.montant_ht;
    document.getElementById('facture-taux-tva').value = facture.taux_tva;
    document.getElementById('facture-description').value = facture.description || '';
    document.getElementById('facture-notes').value = facture.notes || '';
    
    // Modifier le modal pour la mise à jour
    const form = document.getElementById('facture-form');
    form.dataset.editId = id;
    document.getElementById('facture-modal-title').textContent = 'Modifier la facture';
    
    openModal('facture-modal');
}

function generatePDF(id) {
    showNotification('Génération PDF en cours de développement', 'info');
}

// Gestion des modales
function openModal(modalId) {
    document.getElementById(modalId).classList.add('show');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('show');
    
    // Reset form if it's facture modal
    if (modalId === 'facture-modal') {
        const form = document.getElementById('facture-form');
        form.reset();
        delete form.dataset.editId;
        document.getElementById('facture-modal-title').textContent = 'Nouvelle facture';
    }
}

// Administration
async function loadAdminStats() {
    try {
        showLoading();
        const response = await apiCall('/admin/stats');
        const stats = response.data;
        
        document.getElementById('total-clients').textContent = stats.total_clients || 0;
        document.getElementById('total-factures').textContent = stats.total_factures || 0;
        document.getElementById('montant-total').textContent = formatCurrency(stats.montant_total || 0);
        
        // Afficher les factures récentes
        const recentContainer = document.getElementById('recent-factures');
        recentContainer.innerHTML = '';
        
        if (stats.factures_recentes && stats.factures_recentes.length > 0) {
            stats.factures_recentes.forEach(facture => {
                const item = document.createElement('div');
                item.className = 'recent-item';
                item.innerHTML = `
                    <div style="display: flex; justify-content: space-between; align-items: center; padding: 1rem; background: rgba(102, 126, 234, 0.05); border-radius: 10px; margin-bottom: 0.5rem;">
                        <div>
                            <strong>${facture.numero_facture}</strong> - ${facture.client_nom}
                            <br><small>${formatDate(facture.date_facture)}</small>
                        </div>
                        <div style="text-align: right;">
                            <div>${formatCurrency(facture.montant_ttc)}</div>
                            <span class="status ${facture.statut}">${facture.statut}</span>
                        </div>
                    </div>
                `;
                recentContainer.appendChild(item);
            });
        } else {
            recentContainer.innerHTML = '<p>Aucune facture récente</p>';
        }
    } catch (error) {
        showNotification('Erreur lors du chargement des statistiques', 'error');
    } finally {
        hideLoading();
    }
}

// Recherche et filtres
function initFilters() {
    // Filtre pour les factures
    document.getElementById('filter-statut').addEventListener('change', filterFactures);
    document.getElementById('search-facture').addEventListener('input', filterFactures);
    
    // Filtre pour les clients
    document.getElementById('search-client').addEventListener('input', filterClients);
}

function filterFactures() {
    const statutFilter = document.getElementById('filter-statut').value;
    const searchTerm = document.getElementById('search-facture').value.toLowerCase();
    
    const rows = document.querySelectorAll('#factures-table tbody tr');
    
    rows.forEach(row => {
        const statut = row.querySelector('.status').textContent;
        const numero = row.cells[0].textContent.toLowerCase();
        const client = row.cells[1].textContent.toLowerCase();
        
        const matchesStatut = !statutFilter || statut === statutFilter;
        const matchesSearch = !searchTerm || numero.includes(searchTerm) || client.includes(searchTerm);
        
        row.style.display = matchesStatut && matchesSearch ? '' : 'none';
    });
}

function filterClients() {
    const searchTerm = document.getElementById('search-client').value.toLowerCase();
    const rows = document.querySelectorAll('#clients-table tbody tr');
    
    rows.forEach(row => {
        const nom = row.cells[0].textContent.toLowerCase();
        const prenom = row.cells[1].textContent.toLowerCase();
        const entreprise = row.cells[2].textContent.toLowerCase();
        
        const matches = nom.includes(searchTerm) || 
                       prenom.includes(searchTerm) || 
                       entreprise.includes(searchTerm);
        
        row.style.display = matches ? '' : 'none';
    });
}

// Initialisation
document.addEventListener('DOMContentLoaded', function() {
    // Initialiser les onglets
    initTabs();
    
    // Initialiser les filtres
    initFilters();
    
    // Gestionnaires d'événements pour les formulaires
    document.getElementById('client-form').addEventListener('submit', function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        const clientData = Object.fromEntries(formData.entries());
        
        const editId = this.dataset.editId;
        if (editId) {
            updateClient(editId, clientData);
        } else {
            saveClient(clientData);
        }
    });
    
    document.getElementById('facture-form').addEventListener('submit', function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        const factureData = Object.fromEntries(formData.entries());
        
        // Convertir les types appropriés
        factureData.client_id = parseInt(factureData.client_id);
        factureData.montant_ht = parseFloat(factureData.montant_ht);
        factureData.taux_tva = parseFloat(factureData.taux_tva);
        
        const editId = this.dataset.editId;
        if (editId) {
            updateFacture(editId, factureData);
        } else {
            saveFacture(factureData);
        }
    });
    
    // Boutons pour ouvrir les modales
    document.getElementById('nouvelle-facture-btn').addEventListener('click', function() {
        openModal('facture-modal');
    });
    
    document.getElementById('nouveau-client-btn').addEventListener('click', function() {
        document.querySelector('[data-tab="nouveau-client"]').click();
    });
    
    // Bouton pour réinitialiser le formulaire client
    document.getElementById('reset-client-form').addEventListener('click', resetClientForm);
    
    // Gestionnaires pour fermer les modales
    document.querySelectorAll('.modal-close').forEach(btn => {
        btn.addEventListener('click', function() {
            const modal = this.closest('.modal');
            closeModal(modal.id);
        });
    });
    
    // Fermer les modales en cliquant à l'extérieur
    document.querySelectorAll('.modal').forEach(modal => {
        modal.addEventListener('click', function(e) {
            if (e.target === this) {
                closeModal(this.id);
            }
        });
    });
    
    // Charger les données initiales
    loadClients();
    loadFactures();
    
    // Simuler un utilisateur connecté
    currentUser = { username: 'admin', role: 'admin' };
    document.getElementById('username').textContent = currentUser.username;
    
    // Gestionnaire de déconnexion
    document.getElementById('logout-btn').addEventListener('click', function() {
        if (confirm('Êtes-vous sûr de vouloir vous déconnecter ?')) {
            // Ici on pourrait appeler l'API de déconnexion
            showNotification('Déconnexion réussie', 'success');
            // Rediriger vers la page de connexion ou recharger
            location.reload();
        }
    });
    
    // Générer automatiquement le numéro de facture
    document.getElementById('nouvelle-facture-btn').addEventListener('click', function() {
        const year = new Date().getFullYear();
        const count = factures.length + 1;
        const numeroFacture = `FACT-${year}-${count.toString().padStart(4, '0')}`;
        document.getElementById('facture-numero').value = numeroFacture;
        
        // Date par défaut
        const today = new Date().toISOString().split('T')[0];
        document.getElementById('facture-date').value = today;
        
        // Date d'échéance par défaut (30 jours)
        const echeance = new Date();
        echeance.setDate(echeance.getDate() + 30);
        document.getElementById('facture-echeance').value = echeance.toISOString().split('T')[0];
    });
});

