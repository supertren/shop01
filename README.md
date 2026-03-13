# Shop01 — E-Commerce Web Application

An online webshop built with support for multiple payment methods.

## Overview

Shop01 is a full-featured e-commerce website where customers can browse products, add items to a cart, and complete purchases using their preferred payment method.

## Features

- Product catalog with categories and search
- Shopping cart and checkout flow
- User accounts and order history
- Multiple payment method integrations:
  - Credit / Debit Card (Stripe)
  - PayPal
  - Bank Transfer
  - Buy Now Pay Later (e.g. Klarna)
- Order management and confirmation emails
- Admin dashboard for product and order management

## Tech Stack

| Layer    | Technology          |
|----------|---------------------|
| Frontend | HTML / CSS / JS (served by Go) |
| Backend  | Go                  |
| Database | PostgreSQL          |
| Payments | Stripe / PayPal     |

## Getting Started

### Prerequisites

- Go >= 1.22
- PostgreSQL >= 16

### Installation

```bash
git clone https://github.com/your-username/shop01.git
cd shop01
go mod download
```

### Environment Variables

Copy the example env file and fill in your credentials:

```bash
cp .env.example .env
```

Key variables:

```
STRIPE_SECRET_KEY=
PAYPAL_CLIENT_ID=
PAYPAL_CLIENT_SECRET=
KLARNA_CLIENT_ID=
DATABASE_URL=postgres://user:password@localhost:5432/shop01
```

### Database Setup

Make sure PostgreSQL is running:

```bash
sudo systemctl start postgresql
```

Then create the database, set the user password, and apply the schema:

```bash
sudo -u postgres psql -c "CREATE DATABASE shop01;"
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'password';"
sudo -u postgres psql -d shop01 -f schema.sql
```

Update `DATABASE_URL` in your `.env` to match:

```
DATABASE_URL=postgres://postgres:password@localhost:5432/shop01
```

### Sample Data (Optional)

To populate the database with sample products:

```bash
sudo -u postgres psql -d shop01 -c "
INSERT INTO products (name, description, price, image_url, stock) VALUES
('Wireless Headphones', 'High-quality over-ear headphones with noise cancellation.', 79.99, 'https://placehold.co/400x300?text=Headphones', 25),
('Mechanical Keyboard', 'Compact TKL mechanical keyboard with RGB backlight.', 59.99, 'https://placehold.co/400x300?text=Keyboard', 15),
('USB-C Hub', '7-in-1 USB-C hub with HDMI, USB 3.0, and SD card reader.', 34.99, 'https://placehold.co/400x300?text=USB+Hub', 40),
('Webcam 1080p', 'Full HD webcam with built-in microphone and auto-focus.', 49.99, 'https://placehold.co/400x300?text=Webcam', 20);
"
```

### Running Locally

```bash
go run ./cmd/server
```

The server will start at `http://localhost:8080`.

## Project Structure

```
shop01/
├── cmd/
│   └── server/         # Application entry point (main.go)
├── internal/
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Data models
│   ├── payments/       # Payment provider integrations
│   └── store/          # Database access layer
├── web/
│   ├── templates/      # HTML templates
│   └── static/         # CSS, JS, images
├── go.mod
├── go.sum
└── .env.example        # Environment variable template
```

## Payment Integration

Each payment provider is integrated via its official SDK:

- **Stripe** — card payments, webhooks for order confirmation
- **PayPal** — PayPal wallet and card payments via PayPal SDK
- **Bank Transfer** — manual or automated via banking API
- **BNPL** — Klarna or similar, embedded at checkout

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit your changes: `git commit -m "feat: add your feature"`
4. Push and open a Pull Request

## License

MIT
