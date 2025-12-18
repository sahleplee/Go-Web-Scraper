# Go Web Scraper

Bu proje, Go (Golang) ile geliştirilmiş, belirlenen web sitelerinden veri çeken, ekran görüntüsü alan ve linkleri raporlayan güçlü ve hızlı bir araçtır. `chromedp` kütüphanesini kullanarak gerçek bir tarayıcı (headless) üzerinden işlem yapar.

## Özellikler

*   **HTML İndirme**: Hedef sayfanın ham HTML içeriğini kaydeder.
*   **Ekran Görüntüsü**: Sayfanın tam boy ekran görüntüsünü (screenshot) alır.
*   **URL Çıkarma**: Sayfadaki tüm linkleri ayıklar ve listeler.
*   **Eş Zamanlı Tarama**: Birden fazla siteyi aynı anda (concurrent) tarayabilir.
*   **Otomatik Klasörleme**: Çıktıları `html`, `screenshots` ve `url` klasörlerine düzenli bir şekilde kaydeder.
*   **Tarayıcı Desteği**: Google Chrome ve Brave Browser ile uyumludur.

## Gereksinimler

*   [Go](https://go.dev/dl/) (Golang) yüklü olmalıdır.
*   Google Chrome veya Brave Browser yüklü olmalıdır.

## Kurulum

Projeyi bilgisayarınıza indirdikten sonra, gerekli kütüphaneleri yüklemek için terminalde şu komutları çalıştırın:

```bash
# Bağımlılıkları yükle ve güncelle (Önemli: Hataları önlemek için lates sürüm kullanılmalı)
go get -u github.com/chromedp/chromedp@latest github.com/chromedp/cdproto@latest
go mod tidy
```

## Kullanım

Programı terminal veya komut satırı üzerinden çalıştırabilirsiniz.

### 1. Tek Bir Siteyi Tarama
```bash
go run scrapper.go -url="https://www.google.com"
```

### 2. Birden Fazla Siteyi Tarama
Birden fazla URL'yi virgülle ayırarak yazabilirsiniz. Program bunları eş zamanlı olarak tarayacaktır.
```bash
go run scrapper.go -url="https://www.google.com,https://github.com,https://haberler.com"
```

### 3. Brave Browser Kullanarak Tarama
Eğer sisteminizde Chrome yerine Brave yüklü ise `-brave` parametresini ekleyin:
```bash
go run scrapper.go -url="https://www.google.com" -brave
```

### 4. Farklı Bir Tarayıcı Yolu Belirtme
Tarayıcınız standart dışı bir klasörde yüklü ise yolunu (path) belirtebilirsiniz:
```bash
go run scrapper.go -url="https://www.google.com" -exec-path="C:\Program Files\Tarayici\tarayici.exe"
```

## Çıktılar

Program çalıştıktan sonra proje klasöründe şu dizinler oluşur:

*   📂 **html/**: Sitelerin `.html` dosyaları burada saklanır.
*   📂 **screenshots/**: Sitelerin `.png` formatındaki ekran görüntüleri buradadır.
*   📂 **url/**: Her siteden çıkarılan linklerin olduğu `.txt` dosyaları buradadır.

Dosya isimleri taranan sitenin adına göre otomatik oluşturulur (örn: `google.com_screenshot.png`).

## Lisans
Bu proje açık kaynaklıdır ve eğitim amaçlı hazırlanmıştır.
