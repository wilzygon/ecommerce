CREATE TABLE proveedores (
	id UUID NOT NULL,
    cod_proveedor integer,
	nombre VARCHAR(128) NOT NULL,
    ruc_ci VARCHAR(128) NOT NULL,
    telefono VARCHAR(128) NOT NULL,
    direccion VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
	created_at INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM now())::int,
	updated_at INTEGER,
	CONSTRAINT proveedores_id_pk PRIMARY KEY (id)
);

COMMENT ON TABLE proveedores IS 'Almacena los proveedores for the e-commerce';
