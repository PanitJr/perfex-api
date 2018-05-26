INSERT INTO clinic_products(`name`, `related_products`, `price`, `status`,`info`) 
VALUES('โบท็อกแก้ม', '[{"product_id":1,"amount":1}]', 5900, 'Available','รายละเอียดมากมายใครก็ได้บอกที');
INSERT INTO clinic_products(`name`, `related_products`, `price`, `status`,`info`) 
VALUES('โบท็อกคาง', '[{"product_id":1,"amount":1}]', 5500, 'Available','รายละเอียดมากมายใครก็ได้บอกที');
INSERT INTO clinic_products(`name`, `related_products`, `price`, `status`,`info`) 
VALUES('ร้อยไหม', '[{"product_id":2,"amount":1}]', 8800, 'Available','รายละเอียดมากมายใครก็ได้บอกที');
INSERT INTO clinic_products(`name`, `related_products`, `price`, `status`,`info`) 
VALUES('นวดหน้า', '[{"product_id":3,"amount":1}]', 200, 'Available','รายละเอียดมากมายใครก็ได้บอกที');
INSERT INTO clinic_products(`name`, `related_products`, `price`, `status`,`info`) 
VALUES('ขัดผิว', '[{"product_id":4,"amount":1}]', 300, 'Available','รายละเอียดมากมายใครก็ได้บอกที');
INSERT INTO clinic_products(`name`, `related_products`, `price`, `status`,`info`) 
VALUES('เมโส', '[{"product_id":5,"amount":1}]', 2500, 'Available','รายละเอียดมากมายใครก็ได้บอกที');
INSERT INTO clinic_products(`name`, `related_products`, `price`, `status`,`info`) 
VALUES('ฉีดฟิลเลอร์', '[{"product_id":6,"amount":1}]', 9999, 'Available','รายละเอียดมากมายใครก็ได้บอกที');

INSERT INTO products(`name`, `category`, `brand`, `price`,`amount`,`info`,`unit`,`last_update`) 
VALUES('BOTOX A', 'BOTOX', 'American Brand', 100.50,100,'รายละเอียดมากมายใครก็ได้บอกที','u',NOW());
INSERT INTO products(`name`, `category`, `brand`, `price`,`amount`,`info`,`unit`,`last_update`) 
VALUES('ใหมญี่ปุ่น', 'ไหม', 'Japan Brand', 50,100,'รายละเอียดมากมายใครก็ได้บอกที','เส้น',NOW());
INSERT INTO products(`name`, `category`, `brand`, `price`,`amount`,`info`,`unit`,`last_update`) 
VALUES('ชุดนวด', 'อุปกรณ์นวด', 'สามแม่ครัว', 50,100,'รายละเอียดมากมายใครก็ได้บอกที','ชุด',NOW());
INSERT INTO products(`name`, `category`, `brand`, `price`,`amount`,`info`,`unit`,`last_update`) 
VALUES('หินขัด', 'อุปกรณ์ขัด', 'Ital_thai', 5,100,'รายละเอียดมากมายใครก็ได้บอกที','ก้อน',NOW());
INSERT INTO products(`name`, `category`, `brand`, `price`,`amount`,`info`,`unit`,`last_update`) 
VALUES('เมโสโปเตเมีย', 'วัตถุโบราณ', 'ไม่สามารถระบุได้', 100,100,'รายละเอียดมากมายใครก็ได้บอกที','ไม่สามารถนับได้',NOW());
INSERT INTO products(`name`, `category`, `brand`, `price`,`amount`,`info`,`unit`,`last_update`) 
VALUES('ฟิลเลอร์', 'Filler', 'ยาแนวตราตุ๊กแก', 50,100,'รายละเอียดมากมายใครก็ได้บอกที','หลอด',NOW());