const express = require('express');
const cors = require('cors');
const multer = require('multer');
const path = require('path');
const jwt = require('jsonwebtoken');
const fs = require('fs');
const { sequelize, Product, Service, About, Contact, Hero, Footer } = require('./models');
require('dotenv').config();

const app = express();
const port = process.env.PORT || 3000;

// CORS ayarları
app.use(cors({
  origin: [
    'http://localhost:5173',
    'http://admin.localhost:5173',
    'http://klimakozan.com',
    'https://klimakozan.com',
    'http://www.klimakozan.com',
    'https://www.klimakozan.com',
    'http://admin.klimakozan.com',
    'https://admin.klimakozan.com'
  ],
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization']
}));
// Body parser limit ayarı
app.use(express.json({limit: '50mb'}));
app.use(express.urlencoded({limit: '50mb', extended: true}));
app.use('/uploads', express.static(path.join(__dirname, 'uploads')));

// JWT secret key
const JWT_SECRET = process.env.JWT_SECRET;

// Multer yapılandırması
const storage = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, path.join(__dirname, 'uploads'));
  },
  filename: function (req, file, cb) {
    const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9);
    cb(null, uniqueSuffix + path.extname(file.originalname));
  }
});

const upload = multer({
  storage: storage,
  limits: {
    fileSize: 5 * 1024 * 1024 // 5MB limit
  },
  fileFilter: function (req, file, cb) {
    const allowedTypes = /jpeg|jpg|png|gif/;
    const extname = allowedTypes.test(path.extname(file.originalname).toLowerCase());
    const mimetype = allowedTypes.test(file.mimetype);

    if (extname && mimetype) {
      return cb(null, true);
    } else {
      cb('Hata: Sadece resim dosyaları yüklenebilir!');
    }
  }
});

// Auth middleware
const authenticateToken = (req, res, next) => {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1];

  if (token == null) return res.sendStatus(401);

  jwt.verify(token, JWT_SECRET, (err, user) => {
    if (err) return res.sendStatus(403);
    req.user = user;
    next();
  });
};

// Login endpoint
app.post('/api/login', (req, res) => {
  const { username, password } = req.body;
  
  if (username === 'admin' && password === 'admin123') {
    const token = jwt.sign({ username }, JWT_SECRET);
    res.json({ token });
  } else {
    res.status(401).json({ message: 'Geçersiz kullanıcı adı veya şifre' });
  }
});

// Resim yükleme endpoint'i
app.post('/api/admin/upload', authenticateToken, upload.single('file'), (req, res) => {
  try {
    if (!req.file) {
      return res.status(400).json({ message: 'Dosya yüklenemedi' });
    }

    const fileUrl = `http://localhost:${port}/uploads/${req.file.filename}`;
    res.json({ url: fileUrl });
  } catch (error) {
    console.error('Dosya yükleme hatası:', error);
    res.status(500).json({ message: 'Dosya yüklenirken bir hata oluştu' });
  }
});

// Products endpoints
app.get('/api/products', async (req, res) => {
  try {
    const products = await Product.findAll();
    res.json(products);
  } catch (error) {
    res.status(500).json({ message: 'Ürünler yüklenirken bir hata oluştu' });
  }
});

app.get('/api/admin/products', authenticateToken, async (req, res) => {
  try {
    const products = await Product.findAll();
    res.json(products);
  } catch (error) {
    res.status(500).json({ message: 'Ürünler yüklenirken bir hata oluştu' });
  }
});

app.post('/api/admin/products', authenticateToken, async (req, res) => {
  try {
    const product = await Product.create(req.body);
    res.status(201).json(product);
  } catch (error) {
    console.error('Ürün ekleme hatası:', error);
    res.status(500).json({ message: 'Ürün eklenirken bir hata oluştu' });
  }
});

app.put('/api/admin/products/:id', authenticateToken, async (req, res) => {
  try {
    const [updated] = await Product.update(req.body, {
      where: { id: req.params.id }
    });
    
    if (updated) {
      const product = await Product.findByPk(req.params.id);
      res.json(product);
    } else {
      res.status(404).json({ message: 'Ürün bulunamadı' });
    }
  } catch (error) {
    console.error('Ürün güncelleme hatası:', error);
    res.status(500).json({ message: 'Ürün güncellenirken bir hata oluştu' });
  }
});

app.delete('/api/admin/products/:id', authenticateToken, async (req, res) => {
  try {
    const deleted = await Product.destroy({
      where: { id: req.params.id }
    });
    
    if (deleted) {
      res.status(204).send();
    } else {
      res.status(404).json({ message: 'Ürün bulunamadı' });
    }
  } catch (error) {
    console.error('Ürün silme hatası:', error);
    res.status(500).json({ message: 'Ürün silinirken bir hata oluştu' });
  }
});

// Services endpoints
app.get('/api/services', async (req, res) => {
  try {
    const services = await Service.findAll();
    res.json(services);
  } catch (error) {
    res.status(500).json({ message: 'Hizmetler yüklenirken bir hata oluştu' });
  }
});

app.get('/api/admin/services', authenticateToken, async (req, res) => {
  try {
    const services = await Service.findAll();
    res.json(services);
  } catch (error) {
    res.status(500).json({ message: 'Hizmetler yüklenirken bir hata oluştu' });
  }
});

app.post('/api/admin/services', authenticateToken, async (req, res) => {
  try {
    const service = await Service.create(req.body);
    res.status(201).json(service);
  } catch (error) {
    console.error('Hizmet ekleme hatası:', error);
    res.status(500).json({ message: 'Hizmet eklenirken bir hata oluştu' });
  }
});

app.put('/api/admin/services/:id', authenticateToken, async (req, res) => {
  try {
    const [updated] = await Service.update(req.body, {
      where: { id: req.params.id }
    });
    
    if (updated) {
      const service = await Service.findByPk(req.params.id);
      res.json(service);
    } else {
      res.status(404).json({ message: 'Hizmet bulunamadı' });
    }
  } catch (error) {
    console.error('Hizmet güncelleme hatası:', error);
    res.status(500).json({ message: 'Hizmet güncellenirken bir hata oluştu' });
  }
});

app.delete('/api/admin/services/:id', authenticateToken, async (req, res) => {
  try {
    const deleted = await Service.destroy({
      where: { id: req.params.id }
    });
    
    if (deleted) {
      res.status(204).send();
    } else {
      res.status(404).json({ message: 'Hizmet bulunamadı' });
    }
  } catch (error) {
    console.error('Hizmet silme hatası:', error);
    res.status(500).json({ message: 'Hizmet silinirken bir hata oluştu' });
  }
});

// About endpoints
app.get('/api/about', async (req, res) => {
  try {
    const about = await About.findOne();
    res.json(about || {});
  } catch (error) {
    res.status(500).json({ message: 'Hakkımızda bilgileri yüklenirken bir hata oluştu' });
  }
});

app.get('/api/admin/about', authenticateToken, async (req, res) => {
  try {
    const about = await About.findOne();
    res.json(about || {});
  } catch (error) {
    res.status(500).json({ message: 'Hakkımızda bilgileri yüklenirken bir hata oluştu' });
  }
});

app.put('/api/admin/about', authenticateToken, async (req, res) => {
  try {
    let about = await About.findOne();
    
    if (about) {
      about = await about.update(req.body);
    } else {
      about = await About.create(req.body);
    }
    
    res.json(about);
  } catch (error) {
    console.error('Hakkımızda güncelleme hatası:', error);
    res.status(500).json({ message: 'Hakkımızda bilgileri güncellenirken bir hata oluştu' });
  }
});

// Contact endpoints
app.get('/api/contact', async (req, res) => {
  try {
    const contact = await Contact.findOne();
    res.json(contact || {});
  } catch (error) {
    res.status(500).json({ message: 'İletişim bilgileri yüklenirken bir hata oluştu' });
  }
});

app.get('/api/admin/contact', authenticateToken, async (req, res) => {
  try {
    const contact = await Contact.findOne();
    res.json(contact || {});
  } catch (error) {
    res.status(500).json({ message: 'İletişim bilgileri yüklenirken bir hata oluştu' });
  }
});

app.put('/api/admin/contact', authenticateToken, async (req, res) => {
  try {
    let contact = await Contact.findOne();
    
    if (contact) {
      contact = await contact.update(req.body);
    } else {
      contact = await Contact.create(req.body);
    }
    
    res.json(contact);
  } catch (error) {
    console.error('İletişim güncelleme hatası:', error);
    res.status(500).json({ message: 'İletişim bilgileri güncellenirken bir hata oluştu' });
  }
});

// Contact form endpoint
app.post('/api/contact/submit', (req, res) => {
  try {
    console.log('İletişim formu gönderildi:', req.body);
    res.status(200).json({ message: 'Mesajınız başarıyla gönderildi' });
  } catch (error) {
    console.error('Form gönderme hatası:', error);
    res.status(500).json({ message: 'Mesajınız gönderilirken bir hata oluştu' });
  }
});

// Hero section endpoints
app.get('/api/admin/hero', authenticateToken, async (req, res) => {
  try {
    const hero = await Hero.findOne()
    res.json(hero || {
      subheading: 'Klima Kozan\'a Hoş Geldiniz',
      heading: 'Profesyonel Klima ve Beyaz Eşya Çözümleri',
      buttonText: 'Hizmetlerimizi Keşfedin',
      backgroundImage: ''
    })
  } catch (error) {
    res.status(500).json({ error: 'Sunucu hatası' })
  }
})

app.put('/api/admin/hero', authenticateToken, async (req, res) => {
  try {
    const hero = await Hero.findOne()
    if (hero) {
      await Hero.findByIdAndUpdate(hero._id, req.body)
    } else {
      await Hero.create(req.body)
    }
    res.json({ message: 'Hero güncellendi' })
  } catch (error) {
    res.status(500).json({ error: 'Sunucu hatası' })
  }
})

// Footer endpoints
app.get('/api/admin/footer', authenticateToken, async (req, res) => {
  try {
    const footer = await Footer.findOne()
    res.json(footer || {
      copyright: 'Copyright © Klima Kozan 2024',
      socialLinks: {
        twitter: '#',
        facebook: '#',
        instagram: '#'
      },
      links: {
        privacy: 'Gizlilik Politikası',
        terms: 'Kullanım Şartları'
      }
    })
  } catch (error) {
    res.status(500).json({ error: 'Sunucu hatası' })
  }
})

app.put('/api/admin/footer', authenticateToken, async (req, res) => {
  try {
    const footer = await Footer.findOne()
    if (footer) {
      await Footer.findByIdAndUpdate(footer._id, req.body)
    } else {
      await Footer.create(req.body)
    }
    res.json({ message: 'Footer güncellendi' })
  } catch (error) {
    res.status(500).json({ error: 'Sunucu hatası' })
  }
})

// Uploads klasörünü oluştur
const uploadsDir = path.join(__dirname, 'uploads');
if (!fs.existsSync(uploadsDir)) {
  fs.mkdirSync(uploadsDir);
}

// Veritabanı bağlantısı ve sunucuyu başlat
sequelize.sync()
  .then(() => {
    app.listen(port, () => {
      console.log(`Sunucu http://localhost:${port} adresinde çalışıyor`);
    });
  })
  .catch(err => {
    console.error('Veritabanı bağlantı hatası:', err);
  });
