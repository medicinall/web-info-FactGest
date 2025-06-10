-- Ajout de factures de test pour la génération PDF

INSERT INTO factures (numero_facture, client_id, date_facture, date_echeance, statut, montant_ht, taux_tva, montant_tva, montant_ttc, description, notes) VALUES
('FACT-2024-0001', 1, '2024-01-15', '2024-02-15', 'envoyee', 1000.00, 20.0, 200.00, 1200.00, 'Développement site web vitrine', 'Site web responsive avec 5 pages'),
('FACT-2024-0002', 2, '2024-01-20', '2024-02-20', 'payee', 1500.00, 20.0, 300.00, 1800.00, 'Maintenance informatique mensuelle', 'Maintenance préventive et corrective'),
('FACT-2024-0003', 3, '2024-01-25', '2024-02-25', 'brouillon', 800.00, 20.0, 160.00, 960.00, 'Formation bureautique', 'Formation Excel et Word - 8 heures');

