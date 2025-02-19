ALTER TABLE url_mapping 
    ADD CONSTRAINT fk_url_mapping_customer 
    FOREIGN KEY (customer_id) REFERENCES customer (id) 
    ON UPDATE CASCADE ON DELETE RESTRICT;