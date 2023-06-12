CREATE TABLE "products" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "name" varchar UNIQUE NOT NULL,
  "price" int NOT NULL,
  "weight" float NOT NULL,
  "stock" int NOT NULL,
  "image" varchar NOT NULL,
  "description" text NOT NULL,
  "tokopedia_product_id" int NOT NULL,
  "shopee_product_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "total_price" float NOT NULL,
  "tokopedia_order_id" int,
  "shopee_order_id" varchar,
  "customer_name" varchar NOT NULL,
  "customer_phone" varchar NOT NULL,
  "customer_address" varchar NOT NULL,
  "customer_district" varchar NOT NULL,
  "customer_city" varchar NOT NULL,
  "customer_province" varchar NOT NULL,
  "customer_country" varchar NOT NULL,
  "customer_postal_code" varchar NOT NULL,
  "order_status" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "order_products" (
  "order_id" varchar,
  "product_id" varchar,
  "product_name" varchar NOT NULL,
  "product_price" float NOT NULL,
  "product_quantity" int NOT NULL,
  PRIMARY KEY ("order_id", "product_id")
);

CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "shop_name" varchar NOT NULL,
  "tokopedia_fs_id" int NOT NULL,
  "tokopedia_shop_id" int NOT NULL,
  "tokopedia_bearer_token" varchar NOT NULL,
  "shopee_partner_id" int NOT NULL,
  "shopee_shop_id" int NOT NULL,
  "shopee_access_token" varchar NOT NULL,
  "shopee_sign" varchar NOT NULL
);

CREATE INDEX ON "products" ("name");

ALTER TABLE "products" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_products" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
