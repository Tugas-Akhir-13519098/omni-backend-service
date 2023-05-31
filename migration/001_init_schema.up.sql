CREATE TABLE "products" (
  "id" varchar PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "price" int NOT NULL,
  "weight" float NOT NULL,
  "stock" int NOT NULL,
  "image" varchar NOT NULL,
  "description" text NOT NULL,
  "tokopedia_id" int NOT NULL,
  "shopee_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp
);

CREATE INDEX ON "products" ("name");