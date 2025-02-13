alter table url_mapping add constraint short_url_unique unique (short_url); 

alter table url_mapping add constraint original_url_unique unique (original_url); 