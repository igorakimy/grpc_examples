#!/bin/bash

ghz --insecure \
--proto ./ecommerce/product_info.proto \
--call ecommerce.ProductInfo.addProduct \
-d "{\"name\": \"Samsung\", \"description\": \"Samsung Galaxy S10\", \"price\": 700}" \
-n 2000 \
-c 16 \
-O html > output.html \
0.0.0.0:50051