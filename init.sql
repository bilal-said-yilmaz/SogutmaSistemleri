-- Hizmetler
INSERT INTO services (title, description, image) VALUES
('Klima Servisi', 'Kozan''da profesyonel klima bakım, onarım ve montaj hizmetleri.', '/uploads/services/klima-servisi.jpg'),
('Beyaz Eşya Tamiri', 'Buzdolabı, çamaşır makinesi, bulaşık makinesi tamiri ve bakımı.', '/uploads/services/beyaz-esya-tamiri.jpg'),
('Soğutma Sistemleri', 'Endüstriyel soğutma sistemleri kurulum ve bakım hizmetleri.', '/uploads/services/sogutma-sistemleri.jpg');

-- Ürünler
INSERT INTO products (name, description, price, image) VALUES
('Mitsubishi Klima', 'Mitsubishi Electric 12000 BTU Inverter Klima', 24999.99, '/uploads/products/mitsubishi-klima.jpg'),
('Samsung Klima', 'Samsung Wind-Free 18000 BTU Inverter Klima', 29999.99, '/uploads/products/samsung-klima.jpg'),
('LG Klima', 'LG Dual Cool 9000 BTU Inverter Klima', 19999.99, '/uploads/products/lg-klima.jpg');

-- Hakkımızda
INSERT INTO about (title, content, image) VALUES
('Kozan İklimlendirme', 'Kozan''da 15 yılı aşkın süredir klima ve beyaz eşya servisi olarak hizmet vermekteyiz. Profesyonel ekibimiz ve kaliteli hizmet anlayışımızla müşterilerimizin memnuniyetini en üst düzeyde tutmayı hedefliyoruz.', '/uploads/about/about-image.jpg');

-- İletişim
INSERT INTO contact (title, phone, email, address, weekday_hours, saturday_hours, sunday_hours) VALUES
('İletişim Bilgileri', '+90 538 515 3191', 'info@klimakozan.com', 'Kozan, Adana', '09:00 - 18:00', '09:00 - 14:00', 'Kapalı'); 