-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

-- Create ticket_options table
CREATE TABLE IF NOT EXISTS ticket_options (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name character varying,
    "desc" character varying,
    allocation integer,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Create trigger for ticket_options
CREATE TRIGGER update_ticket_options_updated_at
BEFORE UPDATE ON ticket_options
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create purchases table
CREATE TABLE IF NOT EXISTS purchases (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    quantity integer,
    user_id uuid,
    ticket_option_id uuid,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Create trigger for purchases
CREATE TRIGGER update_purchases_updated_at
BEFORE UPDATE ON purchases
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create tickets table
CREATE TABLE IF NOT EXISTS tickets (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    ticket_option_id uuid,
    purchase_id uuid,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Create trigger for tickets
CREATE TRIGGER update_tickets_updated_at
BEFORE UPDATE ON tickets
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Add foreign key constraints
ALTER TABLE purchases
    ADD CONSTRAINT fk_ticket_option
    FOREIGN KEY (ticket_option_id)
    REFERENCES ticket_options(id);

ALTER TABLE tickets
    ADD CONSTRAINT fk_ticket_option
    FOREIGN KEY (ticket_option_id)
    REFERENCES ticket_options(id);

ALTER TABLE tickets
    ADD CONSTRAINT fk_purchase
    FOREIGN KEY (purchase_id)
    REFERENCES purchases(id);

-- Create ar_internal_metadata table (if needed for your ORM)
CREATE TABLE IF NOT EXISTS ar_internal_metadata (
    key character varying NOT NULL,
    value character varying,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (key)
);

-- Create trigger for ar_internal_metadata
CREATE TRIGGER update_ar_internal_metadata_updated_at
BEFORE UPDATE ON ar_internal_metadata
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create schema_migrations table (if needed for your ORM)
CREATE TABLE IF NOT EXISTS schema_migrations (
    version character varying NOT NULL,
    PRIMARY KEY (version)
);
