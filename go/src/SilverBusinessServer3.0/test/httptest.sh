#!/bin/sh

python httptest.py \
AccountTest.test_login_username \
AccountTest.test_register_username \
AccountTest.test_reset_password \
AccountTest.test_get_user_info \
AccountTest.test_get_all_user_info \
AccountTest.test_del_user  \
AccountTest.test_logout_username

