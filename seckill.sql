# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.29)
# Database: seckill
# Generation Time: 2021-01-21 08:23:40 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table order_items
# ------------------------------------------------------------

DROP TABLE IF EXISTS `order_items`;

CREATE TABLE `order_items` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint(20) unsigned NOT NULL,
  `product_id` bigint(20) unsigned NOT NULL,
  `product_sku_id` bigint(20) unsigned NOT NULL,
  `amount` int(10) unsigned NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `rating` int(10) unsigned DEFAULT NULL,
  `review` text COLLATE utf8mb4_unicode_ci,
  `reviewed_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# Dump of table orders
# ------------------------------------------------------------

DROP TABLE IF EXISTS `orders`;

CREATE TABLE `orders` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `no` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `address` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `total_amount` decimal(10,2) NOT NULL COMMENT '订单总金额',
  `remark` text COLLATE utf8mb4_unicode_ci,
  `paid_at` datetime DEFAULT NULL COMMENT '支付时间',
  `payment_method` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '支付方式',
  `payment_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '支付平台订单号',
  `refund_status` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '退款状态',
  `refund_no` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '退款单号',
  `closed` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单是否已关闭',
  `reviewed` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单是否已评价',
  `ship_status` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '物流状态',
  `ship_data` text COLLATE utf8mb4_unicode_ci COMMENT '物流数据',
  `extra` text COLLATE utf8mb4_unicode_ci COMMENT '其他额外数据',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `orders_no_unique` (`no`),
  UNIQUE KEY `orders_refund_no_unique` (`refund_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# Dump of table product_skus
# ------------------------------------------------------------

DROP TABLE IF EXISTS `product_skus`;

CREATE TABLE `product_skus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `stock` int(10) unsigned NOT NULL,
  `product_id` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `product_skus` WRITE;
/*!40000 ALTER TABLE `product_skus` DISABLE KEYS */;

INSERT INTO `product_skus` (`id`, `title`, `description`, `price`, `stock`, `product_id`, `created_at`, `updated_at`)
VALUES
	(1,'iphone12','256G',1.00,0,1,NULL,NULL);

/*!40000 ALTER TABLE `product_skus` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table products
# ------------------------------------------------------------

DROP TABLE IF EXISTS `products`;

CREATE TABLE `products` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `image` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `on_sale` tinyint(1) NOT NULL DEFAULT '1',
  `rating` double(8,2) NOT NULL DEFAULT '5.00',
  `sold_count` int(10) unsigned NOT NULL DEFAULT '0',
  `review_count` int(10) unsigned NOT NULL DEFAULT '0',
  `price` decimal(10,2) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# Dump of table seckill_products
# ------------------------------------------------------------

DROP TABLE IF EXISTS `seckill_products`;

CREATE TABLE `seckill_products` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint(20) unsigned NOT NULL,
  `product_sku_id` bigint(20) unsigned NOT NULL,
  `start_at` datetime NOT NULL,
  `end_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `seckill_products` WRITE;
/*!40000 ALTER TABLE `seckill_products` DISABLE KEYS */;

INSERT INTO `seckill_products` (`id`, `product_id`, `product_sku_id`, `start_at`, `end_at`)
VALUES
	(3,1,1,'2021-01-15 00:00:00','2021-01-17 00:00:00');

/*!40000 ALTER TABLE `seckill_products` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
