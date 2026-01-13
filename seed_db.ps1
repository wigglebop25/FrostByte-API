# seed_db.ps1
# PowerShell script to seed the FrostByte API with initial data.

# Disable SSL Certificate Validation (since we use a self-signed cert)
[System.Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}

$BaseUrl = "https://frostbyte-api.southeastasia.cloudapp.azure.com/api/v1"

# 1. Ask for the Admin JWT Token
$Token = Read-Host "Please enter your Admin JWT Token (Bearer token)"
$Headers = @{
    "Content-Type"  = "application/json"
    "Authorization" = "Bearer $Token"
}

# --- Data Definition ---

$Categories = @(
    @{ name = "Zensai"; description = "Appetizers and Starters" },
    @{ name = "Sushi & Sashimi"; description = "Fresh Sushi and Sashimi" },
    @{ name = "Menrui"; description = "Noodle Dishes" },
    @{ name = "Donburi"; description = "Rice Bowls" },
    @{ name = "Kanmi"; description = "Desserts and Sweets" }
)

$Products = @(
    @{
        name = "Pork Gyoza"
        description = "Six pieces of pan-seared dumplings served with a soy-vinegar dipping sauce"
        price = 8.50
        product_image_uri = "/images/pork-gyoza.jpg"
        categories = @("Zensai")
    },
    @{
        name = "Edamame"
        description = "Steamed young soybeans tossed in flakey sea salt"
        price = 5.00
        product_image_uri = "/images/edamame.jpg"
        categories = @("Zensai")
    },
    @{
        name = "Shrimp Tempura"
        description = "Three pieces of lightly battered, crispy fried shrimp"
        price = 11.00
        product_image_uri = "/images/shrimp-tempura.jpg"
        categories = @("Zensai")
    },
    @{
        name = "Maguro Nigiri"
        description = "Two pieces of fresh Bluefin tuna over hand-pressed vinegared rice"
        price = 12.00
        product_image_uri = "/images/maguro-nigiri.jpg"
        categories = @("Sushi & Sashimi")
    },
    @{
        name = "California Roll"
        description = "Classic roll with crab mix, avocado, and cucumber"
        price = 9.50
        product_image_uri = "/images/cali-roll.jpg"
        categories = @("Sushi & Sashimi")
    },
    @{
        name = "Salmon Sashimi"
        description = "Five thick slices of premium fresh Atlantic salmon"
        price = 14.50
        product_image_uri = "/images/salmon-sashimi.jpg"
        categories = @("Sushi & Sashimi")
    },
    @{
        name = "Tonkotsu Ramen"
        description = "Rich pork bone broth, chashu pork, bamboo shoots, and a soft-boiled egg"
        price = 15.99
        product_image_uri = "/images/tonkotsu-ramen.jpg"
        categories = @("Menrui")
    },
    @{
        name = "Tempura Udon"
        description = "Thick wheat noodles in a clear dashi broth topped with shrimp tempura"
        price = 13.50
        product_image_uri = "/images/tempura-udon.jpg"
        categories = @("Menrui")
    },
    @{
        name = "Vegetable Yakisoba"
        description = "Stir-fried buckwheat noodles with cabbage, carrots, and savory sauce"
        price = 12.00
        product_image_uri = "/images/yakisoba.jpg"
        categories = @("Menrui")
    },
    @{
        name = "Gyu-Don"
        description = "Thinly sliced beef and onions simmered in a sweet soy dashi over rice"
        price = 12.50
        product_image_uri = "/images/gyudon.jpg"
        categories = @("Donburi")
    },
    @{
        name = "Katsu-Don"
        description = "Crispy pork cutlet and egg simmered in savory broth over rice"
        price = 13.50
        product_image_uri = "/images/katsudon.jpg"
        categories = @("Donburi")
    },
    @{
        name = "Unagi-Don"
        description = "Grilled freshwater eel glazed with sweet tare sauce over steamed rice"
        price = 21.00
        product_image_uri = "/images/unagi-don.jpg"
        categories = @("Donburi")
    },
    @{
        name = "Matcha Mochi"
        description = "Sweet glutinous rice cake filled with premium green tea cream"
        price = 6.50
        product_image_uri = "/images/matcha-mochi.jpg"
        categories = @("Kanmi")
    },
    @{
        name = "Taiyaki"
        description = "Fish-shaped waffle cake filled with sweet red bean paste"
        price = 7.00
        product_image_uri = "/images/taiyaki.jpg"
        categories = @("Kanmi")
    },
    @{
        name = "Black Sesame Ice Cream"
        description = "Creamy, nutty, and slightly savory roasted black sesame frozen treat"
        price = 5.50
        product_image_uri = "/images/sesame-ice-cream.jpg"
        categories = @("Kanmi")
    }
)

# --- Execution ---

Write-Host "`n--- Creating Categories ---" -ForegroundColor Cyan

foreach ($cat in $Categories) {
    try {
        $body = $cat | ConvertTo-Json
        $response = Invoke-RestMethod -Uri "$BaseUrl/categories" -Method Post -Headers $Headers -Body $body
        Write-Host "Created Category: $($cat.name)" -ForegroundColor Green
    }
    catch {
        Write-Host "Failed to create category $($cat.name): $_" -ForegroundColor Red
    }
}

Write-Host "`n--- Creating Products ---" -ForegroundColor Cyan

foreach ($prod in $Products) {
    try {
        $body = $prod | ConvertTo-Json -Depth 5
        $response = Invoke-RestMethod -Uri "$BaseUrl/products" -Method Post -Headers $Headers -Body $body
        Write-Host "Created Product: $($prod.name)" -ForegroundColor Green
    }
    catch {
        Write-Host "Failed to create product $($prod.name): $_" -ForegroundColor Red
    }
}

Write-Host "`nSeed Complete!" -ForegroundColor Cyan
