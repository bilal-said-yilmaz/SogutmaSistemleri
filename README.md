# Proje Adı

Bu proje, Vue.js, Go, PostgreSQL ve Docker kullanılarak geliştirilmiş bir web sitesidir. Projede, dinamik içerik yönetimi yapılabilen bir **admin paneli** ve **kullanıcı dostu bir ana sayfa** bulunmaktadır.

## Proje Genel Özellikleri

- **Tek Sayfa (SPA) Yapısı**  
  - Kullanıcılar, sayfanın üst kısmındaki menüden aşağıdaki bölülere erişebilecektir:
    - **Hizmetlerimiz**
    - **Ürünlerimiz**
    - **Hakkımızda**
    - **İletişim**
    - **Admin Paneline Giriş**
  - **Sayfa Yenilenmeden Geçiş:** Seçeneklere tıklandığında sayfa yenilenmeden aşağı kayarak ilgili bölüme geçiş sağlanacaktır.

- **Admin Paneli**  
  - Yetkili kullanıcılar, **admin paneli** aracılığıyla site içeriğini yönetebilir.
  - **Ürünler**, **hizmetler**, **hakkımızda** ve **iletişim** bölülerinin içeriğini ekleyebilir, düzenleyebilir veya silebilir.

- **Veritabanı**  
  - PostgreSQL kullanılacak ve dinamik veriler burada saklanacaktır.

- **Docker Desteği**  
  - Proje, Docker kullanılarak kolayca kurulabilir hale getirilecektir.

## Kullanılan Teknolojiler

- **Frontend:** Vue.js
- **Backend:** Go
- **Veritabanı:** PostgreSQL
- **Containerization:** Docker
- **Paket Yöneticisi:** Yarn

## Kurulum

1. **Depoyu Klonlayın:**  
   ```sh
   git clone <repo-link>
   cd <proje-klasörü>
   ```

2. **Docker ile Başlatmak:**  
   ```sh
   docker-compose up --build
   ```

3. **Manuel Kurulum:**  
   - PostgreSQL'i kurun ve veritabanı bağlantı bilgilerini ayarlayın.
   - Backend'i derleyin ve çalıştırın:
     ```sh
     go run main.go
     ```
   - Vue.js projesini başlatın:
     ```sh
     cd frontend
     yarn install
     yarn dev
     ```

## Katkıda Bulunma
Katkıda bulunmak için pull request açabilirsiniz. Sorular ve öneriler için issue oluşturabilirsiniz.

---
**Lisans:** MIT 