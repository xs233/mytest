#!/bin/sh

python httptest.py \
AccountTest.test_login_username \
FavoriteTest.test_save_seach_person_rule \
FavoriteTest.test_query_seach_person_rule \
FavoriteTest.test_query_favorite_device \
FavoriteTest.test_add_favorite \
FavoriteTest.test_query_favorite_by_page  \
FavoriteTest.test_cancel_favorite 
