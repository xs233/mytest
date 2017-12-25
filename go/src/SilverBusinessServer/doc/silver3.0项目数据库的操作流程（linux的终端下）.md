silver3.0项目数据库的操作流程（linux的终端下）

-------------------------------------------------------------------------------------------

mysql -h 192.168.0.200 -P 3306 -u impower -p

密码:impower2017

进入服务器的mysql后的操作是

将db_setup.sql文件中的命令全部复制拷贝到终端中，回车键



-----------

下面是创建数据的操作

use impowerdb;

SHOW TABLES;

SHOW COLUMNS FROM imp_t_device;

INSERT INTO `imp_t_device` VALUES (1, 'vms_id', 'camera', '192.168.0.1', 'rtsp_url','main_stream_url','sub_stream_url','access_account','111111');

SHOW INDEX FROM imp_t_device;



update imp_t_user set password = '96E79218965EB72C92A549DD5A330112';

select password from imp_t_user;