-- MySQL dump 10.13  Distrib 8.4.7, for Linux (x86_64)
--
-- Host: localhost    Database: fiber_learning
-- ------------------------------------------------------
-- Server version	8.4.7

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `annotations`
--

DROP TABLE IF EXISTS `annotations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `annotations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `content` text NOT NULL,
  `post_id` bigint unsigned DEFAULT NULL,
  `line_number` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_annotations_deleted_at` (`deleted_at`),
  KEY `idx_annotations_post_id` (`post_id`),
  KEY `idx_annotations_line_number` (`line_number`),
  CONSTRAINT `fk_annotations_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=82 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `annotations`
--

LOCK TABLES `annotations` WRITE;
/*!40000 ALTER TABLE `annotations` DISABLE KEYS */;
INSERT INTO `annotations` VALUES (3,'2025-10-29 17:46:30.577','2025-10-29 17:46:30.577',NULL,'Cài Kubespray (clone từ git) -> control-plane',8,4),(4,'2025-10-29 17:46:48.975','2025-10-29 17:46:48.975',NULL,'nếu có rồi thì bỏ qua',8,8),(5,'2025-10-29 17:47:18.288','2025-10-29 17:47:18.288',NULL,'có thể copy sẵn file template sample của Kubespray: cp -r inventory/sample inventory/mycluster',8,12),(6,'2025-10-29 17:47:37.833','2025-10-29 17:47:37.833',NULL,'nếu chưa có nano thì chạy command: apt update && install nano',8,13),(7,'2025-10-29 17:48:26.719','2025-10-29 17:48:26.719','2025-10-29 17:54:36.247','Thay bằng các IP thực tế',8,15),(8,'2025-10-29 17:51:46.642','2025-10-29 17:51:46.642','2025-10-29 17:51:58.181','ansible_host: IP hoặc hostname mà Ansible sẽ kết nối đến.\n\nip: IP nội bộ của node trong cluster Kubernetes (dùng trong kubeadm).\n\nansible_user: user để Ansible SSH vào node.',8,16),(9,'2025-10-29 17:51:58.185','2025-10-29 17:51:58.185','2025-10-29 17:52:06.677','- ansible_host: IP hoặc hostname mà Ansible sẽ kết nối đến.\n\n- ip: IP nội bộ của node trong cluster Kubernetes (dùng trong kubeadm).\n\n- ansible_user: user để Ansible SSH vào node.',8,16),(10,'2025-10-29 17:52:06.684','2025-10-29 17:52:06.684','2025-10-29 17:52:30.744','- ansible_host: IP hoặc hostname mà Ansible sẽ kết nối đến.\n\n\n- ip: IP nội bộ của node trong cluster Kubernetes (dùng trong kubeadm).\n\n- ansible_user: user để Ansible SSH vào node.',8,16),(11,'2025-10-29 17:52:30.751','2025-10-29 17:52:30.751',NULL,'- ansible_host: IP hoặc hostname mà Ansible sẽ kết nối đến.\n// ip: IP nội bộ của node trong cluster Kubernetes (dùng trong kubeadm).\n// ansible_user: user để Ansible SSH vào node.',8,16),(12,'2025-10-29 17:54:36.253','2025-10-29 17:54:36.253',NULL,'danh sách tất cả các node trong cluster Kubernetes',8,15),(13,'2025-10-29 17:55:09.227','2025-10-29 17:55:09.227',NULL,'node chạy control-plane (Master nodes)',8,22),(14,'2025-10-29 17:55:43.988','2025-10-29 17:55:43.988',NULL,'Control-plane node chạy:\n\n// kube-apiserver\n\n// kube-scheduler\n\n// kube-controller-manager\n\n// Chỉ định node nào làm master để Kubespray triển khai.',8,23),(15,'2025-10-29 17:57:09.196','2025-10-29 17:57:09.196','2025-10-29 17:57:25.928','Liệt kê các node chạy etcd (cơ sở dữ liệu của Kubernetes)',8,26),(16,'2025-10-29 17:57:25.935','2025-10-29 17:57:25.935',NULL,'Liệt kê các node chạy etcd (cơ sở dữ liệu của Kubernetes)\n//chọn các control-plane node làm etcd để cluster gọn và ổn định',8,26),(17,'2025-10-29 17:57:39.426','2025-10-29 17:57:39.426','2025-10-29 17:58:03.199','các worker nodes trong cluster',8,30),(18,'2025-10-29 17:58:03.206','2025-10-29 17:58:03.206','2025-10-30 15:33:35.390','các worker nodes trong cluster\n// Chúng sẽ chạy workloads/pods, nhưng không chạy các thành phần control-plane.\n// Kubespray sẽ cài kubelet, kube-proxy, containerd… lên các node này',8,30),(19,'2025-10-29 17:58:39.977','2025-10-29 17:58:39.977','2025-10-29 17:58:53.655','Đây là group chứa tất cả node của cluster Kubernetes',8,34),(20,'2025-10-29 17:58:53.662','2025-10-29 17:58:53.662','2025-10-30 15:33:49.718','Đây là group chứa tất cả node của cluster Kubernetes //children nghĩa là nhóm con, tổng hợp node control-plane + node worker',8,34),(21,'2025-10-29 17:59:40.691','2025-10-29 17:59:40.691',NULL,'ho phép Ansible chạy các lệnh với quyền root (sudo)',8,38),(22,'2025-10-29 17:59:51.831','2025-10-29 17:59:51.831',NULL,'Chỉ định private key để Ansible SSH đến các node mà không cần nhập password.',8,39),(23,'2025-10-29 18:02:31.231','2025-10-29 18:02:31.231',NULL,'chạy trên node control_plane có file admin.conf',8,54),(24,'2025-10-29 18:04:01.587','2025-10-29 18:04:01.587','2025-10-30 14:40:17.927','Ansible là công cụ bắt buộc để Kubespray hoạt động, chỉ cài trên máy proxy/control-plane để điều khiển cluster.',8,1),(25,'2025-10-29 18:27:50.637','2025-10-29 18:27:50.637','2025-10-29 18:28:00.398','Vào môi trường ảo của python',8,10),(26,'2025-10-29 18:28:00.407','2025-10-29 18:28:00.407',NULL,'Vào môi trường ảo của python trên linux server',8,10),(27,'2025-10-30 14:40:17.935','2025-10-30 14:40:17.935','2025-10-30 14:40:35.719','Ansible là công cụ bắt buộc để\n Kubespray hoạt động, chỉ cài trên máy proxy/control-plane để điều khiển cluster.',8,1),(28,'2025-10-30 14:40:35.725','2025-10-30 14:40:35.725','2025-10-30 14:40:45.311','Ansible là công cụ bắt buộc để cho\n Kubespray hoạt động, chỉ cài trên máy proxy/control-plane để điều khiển cluster.',8,1),(29,'2025-10-30 14:40:45.319','2025-10-30 14:40:45.319','2025-10-30 14:40:52.646','Ansible là công cụ bắt buộc để   \n Kubespray hoạt động, chỉ cài trên máy proxy/control-plane để điều khiển cluster.',8,1),(30,'2025-10-30 14:40:52.654','2025-10-30 14:40:52.654','2025-10-30 14:41:03.966','Ansible là công cụ bắt buộc để     Kubespray hoạt động, chỉ cài trên máy proxy/control-plane để điều khiển cluster.',8,1),(31,'2025-10-30 14:41:03.973','2025-10-30 14:41:03.973','2025-10-30 14:41:22.116','Ansible là công cụ bắt buộc cho Kubespray hoạt động, chỉ cài trên máy proxy/control-plane để điều khiển cluster.',8,1),(32,'2025-10-30 14:41:22.122','2025-10-30 14:41:22.122',NULL,'Ansible là công cụ bắt buộc cho Kubespray hoạt động, chỉ cài trên máy proxy/control-plane để có thể điều khiển cluster.',8,1),(33,'2025-10-30 15:31:17.810','2025-10-30 15:31:17.810','2025-10-31 08:48:40.848','123',8,2),(34,'2025-10-30 15:33:35.395','2025-10-30 15:33:35.395',NULL,'',8,30),(35,'2025-10-30 15:33:39.111','2025-10-30 15:33:39.111',NULL,'các worker nodes trong cluster\n// Chúng sẽ chạy workloads/pods, nhưng không chạy các thành phần control-plane.\n// Kubespray sẽ cài kubelet, kube-proxy, containerd… lên các node này',8,31),(36,'2025-10-30 15:33:49.726','2025-10-30 15:33:49.726',NULL,'',8,34),(37,'2025-10-30 15:33:55.252','2025-10-30 15:33:55.252',NULL,'Đây là group chứa tất cả node của cluster Kubernetes //children nghĩa là nhóm con, tổng hợp node control-plane + node worker',8,35),(38,'2025-10-31 08:48:40.855','2025-10-31 08:48:40.855',NULL,'',8,2),(39,'2025-10-31 09:28:32.119','2025-10-31 09:28:32.119',NULL,'Hiển thị tất cả tập tin, bao gồm các tập tin ẩn',9,5),(40,'2025-10-31 09:28:47.166','2025-10-31 09:28:47.166',NULL,'Liệt kê kèm thông tin mô tả chi tiết',9,4),(41,'2025-10-31 09:29:57.381','2025-10-31 09:29:57.381',NULL,'Hiển thị tất cả tập tin (không bao gồm tập tin ẩn)',9,3),(42,'2025-10-31 09:31:09.942','2025-10-31 09:31:09.942',NULL,'Đi đến thư mục được chỉ định',9,7),(43,'2025-10-31 09:31:20.414','2025-10-31 09:31:20.414',NULL,'Quay lại thư mục trước đó',9,8),(44,'2025-10-31 09:31:26.549','2025-10-31 09:31:26.549',NULL,'Đi đến thư mục chính (home)',9,9),(45,'2025-10-31 09:31:55.909','2025-10-31 09:31:55.909',NULL,'Tạo thư mục mới',9,11),(46,'2025-10-31 09:32:38.582','2025-10-31 09:32:38.582',NULL,'Xóa thư mục được chỉ định',9,13),(47,'2025-10-31 09:33:07.869','2025-10-31 09:33:07.869','2025-10-31 09:33:14.769','# Sao chép thư mục một cách đệ quy (recursive) (gồm cả file và thư mục con)',9,16),(48,'2025-10-31 09:33:14.774','2025-10-31 09:33:14.774',NULL,'Sao chép thư mục một cách đệ quy (recursive) (gồm cả file và thư mục con)',9,16),(49,'2025-10-31 09:33:40.941','2025-10-31 09:33:40.941',NULL,'Sao chép tập tin (không phải thư mục)',9,15),(50,'2025-10-31 09:34:00.237','2025-10-31 09:34:00.237',NULL,'Đổi tên thư mục',9,18),(51,'2025-10-31 09:35:00.026','2025-10-31 09:35:00.026',NULL,'Di chuyển tập tin đến thư mục đích',9,19),(52,'2025-10-31 09:35:21.838','2025-10-31 09:35:21.838',NULL,'Xóa tập tin đơn (không xóa được thư mục lớn)',9,21),(53,'2025-10-31 09:35:53.742','2025-10-31 09:35:53.742',NULL,'Xóa thư mục một cách đệ quy (bao gồm tất cả file và folder con)',9,22),(54,'2025-10-31 09:36:40.742','2025-10-31 09:36:40.742','2025-10-31 09:36:46.186','Hiển thị 10 dòng đầu tiên',9,31),(55,'2025-10-31 09:36:46.194','2025-10-31 09:36:46.194',NULL,'Hiển thị 10 dòng đầu tiên của tập tin',9,31),(56,'2025-10-31 09:36:58.533','2025-10-31 09:36:58.533',NULL,'Hiển thị 5 dòng đầu tiên của tập tin',9,32),(57,'2025-10-31 09:37:15.262','2025-10-31 09:37:15.262',NULL,'Hiển thị 10 dòng cuối cùng của tập tin',9,34),(58,'2025-10-31 09:37:27.414','2025-10-31 09:37:27.414',NULL,'Hiển thị 5 dòng cuối cùng của tập tin',9,35),(59,'2025-10-31 09:38:07.470','2025-10-31 09:38:07.470',NULL,'Tìm kiếm dựa vào tên tập tin',9,37),(60,'2025-10-31 09:38:17.574','2025-10-31 09:38:17.574','2025-10-31 09:38:22.018','# Tìm kiếm đệ quy trong thư mục',9,38),(61,'2025-10-31 09:38:22.025','2025-10-31 09:38:22.025',NULL,'Tìm kiếm đệ quy trong thư mục',9,38),(62,'2025-10-31 09:38:39.853','2025-10-31 09:38:39.853',NULL,'Gán quyền đọc, ghi, thực thi cho chủ sở hữu và quyền đọc, thực thi cho người khác',9,41),(63,'2025-10-31 09:38:47.732','2025-10-31 09:38:47.732',NULL,'Gán quyền thực thi cho tập tin',9,42),(64,'2025-10-31 09:39:28.885','2025-10-31 09:39:28.885',NULL,'Đặt quyền mặc định là 755 cho thư mục và 644 cho tập tin',9,46),(65,'2025-10-31 09:39:41.637','2025-10-31 09:39:41.637',NULL,'Hiển thị tất cả các tiến trình',9,50),(66,'2025-10-31 09:41:36.899','2025-10-31 09:41:36.899',NULL,'yêu cầu tiến trình tự dừng một cách an toàn, cho phép nó giải phóng tài nguyên, lưu trạng thái, đóng file…',9,54),(67,'2025-10-31 09:41:58.290','2025-10-31 09:41:58.290',NULL,'Gửi tín hiệu SIGKILL (9), bắt buộc dừng ngay lập tức.\n\nTiến trình không có cơ hội cleanup, không thể trap tín hiệu này.\n\nDùng khi tiến trình treo, không phản hồi với SIGTERM.',9,55),(68,'2025-10-31 09:42:27.453','2025-10-31 09:42:27.453',NULL,'Định dạng dễ đọc cho con người',9,60),(69,'2025-10-31 09:42:45.964','2025-10-31 09:42:45.964',NULL,'Định dạng dễ đọc',9,64),(70,'2025-10-31 09:43:00.550','2025-10-31 09:43:00.550',NULL,'Hiển thị toàn bộ thông tin hệ thống',9,66),(71,'2025-10-31 09:43:19.962','2025-10-31 09:43:19.962',NULL,'Hiển thị địa chỉ IP của các giao diện mạng',9,81),(72,'2025-10-31 09:43:26.621','2025-10-31 09:43:26.621',NULL,'Hiển thị bảng định tuyến',9,82),(73,'2025-10-31 09:44:20.422','2025-10-31 09:44:20.422',NULL,'Cập nhật danh sách package',9,89),(74,'2025-10-31 09:44:28.050','2025-10-31 09:44:28.050',NULL,'Cài đặt một package',9,90),(75,'2025-10-31 09:44:32.837','2025-10-31 09:44:32.837',NULL,'Gỡ bỏ một package',9,91),(76,'2025-10-31 09:44:41.939','2025-10-31 09:44:41.939',NULL,'Cập nhật danh sách package',9,93),(77,'2025-10-31 09:44:50.674','2025-10-31 09:44:50.674',NULL,'Cài đặt một package',9,94),(78,'2025-10-31 09:45:06.946','2025-10-31 09:45:06.946',NULL,'Gỡ bỏ một package',9,95),(79,'2025-10-31 09:47:27.837','2025-10-31 09:47:27.837',NULL,'Tạo một lưu trữ nén',9,98),(80,'2025-10-31 09:47:34.786','2025-10-31 09:47:34.786',NULL,'Giải nén một lưu trữ nén',9,99),(81,'2025-10-31 09:47:54.829','2025-10-31 09:47:54.829',NULL,'Tạo phím tắt cho \'ls -la\'',9,110);
/*!40000 ALTER TABLE `annotations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `content` text NOT NULL,
  `post_id` bigint unsigned DEFAULT NULL,
  `author_id` bigint unsigned DEFAULT NULL,
  `line_number` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_comments_deleted_at` (`deleted_at`),
  KEY `idx_comments_post_id` (`post_id`),
  KEY `idx_comments_author_id` (`author_id`),
  KEY `idx_comments_line_number` (`line_number`),
  CONSTRAINT `fk_posts_comments` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_users_comments` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (4,'2025-10-30 10:05:48.092','2025-10-30 10:05:48.092',NULL,'comment',8,2,1);
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `summary` varchar(255) DEFAULT NULL,
  `content` text,
  `author_id` bigint unsigned DEFAULT NULL,
  `cover_url` varchar(512) DEFAULT NULL,
  `tags` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_posts_deleted_at` (`deleted_at`),
  KEY `idx_posts_author_id` (`author_id`),
  CONSTRAINT `fk_users_posts` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` VALUES (8,'2025-10-29 17:29:35.040','2025-10-31 08:36:13.935',NULL,'Cài đặt k8s trên các nodes (từ proxy)','Cài k8s trên proxy cho các host con','<p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">1. Cài ansible</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">sudo apt update</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">sudo apt install -y ansible sshpass python3</span></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">Cài Kubespray (clone từ git) -&gt; control-plane</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">git clone </span><a href=\"https://github.com/kubernetes-sigs/kubespray.git\" rel=\"noopener noreferrer\" target=\"_blank\" style=\"background-color: transparent; color: rgb(255, 255, 255);\">https://github.com/kubernetes-sigs/kubespray.git</a></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">cd kubespray</span></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">-&gt; cài các dependencies</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">sudo apt update</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">sudo apt install -y python3-pip (nếu có rồi thì bỏ qua)</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">python3 -m venv venv</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">source venv/bin/activate</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">pip install -r requirements.txt</span></p><p><br></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">2.Tạo file ini (hoặc có thể copy sẵn file template sample của Kubespray: cp -r inventory/sample inventory/mycluster)</strong></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">nano inventory.ini</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">(nếu chưa có nano thì chạy command: apt update &amp;&amp; install nano)</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">*nội dung file*</span></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">[all]</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node1 ansible_host=192.168.1.51 ip=192.168.1.51 ansible_user=dev</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node2 ansible_host=192.168.1.52 ip=192.168.1.52 ansible_user=dev</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node3 ansible_host=192.168.1.68 ip=192.168.1.68 ansible_user=dev</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node4 ansible_host=192.168.1.73 ip=192.168.1.73 ansible_user=dev</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node5 ansible_host=192.168.1.109 ip=192.168.1.109 ansible_user=dev</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node6 ansible_host=192.168.1.110 ip=192.168.1.110 ansible_user=dev</span></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">[kube_control_plane]</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node1</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node2</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node3</span></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">[etcd]</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node1</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node2</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node3</span></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">[kube_node]</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node4</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node5</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">node6</span></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">[k8s_cluster:children]</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">kube_control_plane</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">kube_node</span></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">[k8s_cluster:vars]</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">ansible_become=true</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">ansible_ssh_private_key_file=~/.ssh/id_rsa</span></p><p><br></p><p><br></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">*****</span></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">Xong rồi tạo file key:</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">ssh-keygen -t rsa -b 4096 -C \"ansible-key\"</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">*****</span></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">3. ping xem kết nối (tùy chọn)</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">ansible all -i inventory.ini -m ping</span></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">Copy các key qua các nodes:</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">ssh-copy-id dev@192.168.1...</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">ssh-copy-id dev@192.168.1...</span></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">4.Cài đặt Kubernetes (containerd + kubeadm/kubectl/kubelet)</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">cp -rfp inventory/sample inventory/mycluster</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">mv ~/inventory.ini inventory/mycluster/hosts.ini</span></p><p><br></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">5. Chạy playbook cài đặt</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">ansible-playbook -i inventory/mycluster/hosts.ini cluster.yml -b -v --ask-become-pass</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">&nbsp;(đợi sẽ khá lâu và sẽ hỏi password của các servers (nhập 1 lần) - nhớ thay đường dẫn bằng đường dẫn chính xác đến file .ini)</span></p><p><br></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">6. Cấu hình kubectl (chạy trên node control_plane có file admin.conf)</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">mkdir -p $HOME/.kube</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">sudo chown $(id -u):$(id -g) $HOME/.kube/config</span></p><p><br></p><p><strong style=\"background-color: transparent; color: rgb(255, 255, 255);\">7. Kiểm tra cluster</strong></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">kubectl get nodes -o wide</span></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">kubectl get pods -A</span></p><p><br></p><p><span style=\"background-color: transparent; color: rgb(255, 255, 255);\">Khi tất cả node có STATUS = Ready → cluster đã hoạt động hoàn chỉnh</span></p><p><br></p>',2,'https://kb.pavietnam.vn/wp-content/uploads/2021/08/k8s-logo-1200x900.png','k8s'),(9,'2025-10-31 09:01:44.199','2025-10-31 09:49:05.309',NULL,'TỔNG HỢP NHỮNG CÂU LỆNH LINUX THÔNG DỤNG','TỔNG HỢP NHỮNG CÂU LỆNH LINUX THÔNG DỤNG','<h2>Thao Tác Tập Tin và Thư Mục</h2><p><strong>ls – Liệt kê nội dung của một thư mục.</strong></p><p>ls</p><p>ls -l&nbsp;&nbsp;&nbsp;</p><p>ls -a&nbsp;</p><p><br></p><p><strong style=\"color: rgb(56, 189, 248);\">cd – Thay đổi thư mục hiện tại.</strong></p><p>cd /đường/dẫn/đến/thư/mục</p><p>cd ..&nbsp;&nbsp;&nbsp;&nbsp;</p><p>cd ~&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</p><h3><br></h3><p><strong>mkdir – Tạo một thư mục mới.</strong></p><p>mkdir thư_mục_mới</p><p><br></p><p><strong>rmdir – Xóa bỏ thư mục.</strong></p><p>rmdir tên_thư_mục</p><p><br></p><p><strong>cp – Sao chép tập tin hoặc thư mục.</strong></p><p>cp tập_tin_nguồn đích</p><p>cp -r thư_mục_nguồn thư_mục_đích&nbsp;</p><p><br></p><p><strong>mv – Di chuyển hoặc đổi tên tập tin và thư mục.</strong></p><p>mv tên_cũ tên_mới</p><p>mv tên_tập_tin /đường/dẫn/đến/đích/</p><p><br></p><p><strong>rm – Xóa tập tin hoặc thư mục.</strong></p><p>rm tên_tập_tin</p><p>rm -r tên_thư_mục&nbsp;</p><p><br></p><p><strong>touch – Tạo một tập tin trống hoặc cập nhật dấu thời gian của tập tin hiện có.</strong></p><p>touch tên_tập_tin</p><p><br></p><h2><strong style=\"color: rgb(255, 255, 255);\">Xem &amp; Thao Tác Tập Tin</strong></h2><p><strong>cat – Hiển thị nội dung của một tập tin.</strong></p><p>cat tên_tập_tin</p><p><br></p><p><strong>less – Xem nội dung tập tin (với chế độ xem từng trang).</strong></p><p>less tên_tập_tin</p><p><br></p><p><strong>head – Hiển thị 10 dòng đầu tiên của một tập tin (mặc định).</strong></p><p>head tên_tập_tin</p><p>head -n 5 tên_tập_tin&nbsp;</p><p><br></p><p><strong>tail – Hiển thị 10 dòng cuối cùng của một tập tin (mặc định).</strong></p><p>tail tên_tập_tin</p><p>tail -n 5 tên_tập_tin</p><p><br></p><p><strong>grep – Tìm kiếm mẫu trong các tập tin. (searching, filtering)</strong></p><p>grep \'cụm_từ_tìm_kiếm\' tên_tập_tin</p><p>grep -r \'cụm_từ_tìm_kiếm\' /đường/dẫn/đến/thư/mục&nbsp;&nbsp;</p><p><br></p><h2>Quyền &amp; Sở Hữu</h2><p><strong>chmod – Thay đổi quyền của tập tin.</strong></p><p>chmod 755 tên_tập_tin&nbsp;</p><p>chmod +x script.sh&nbsp;</p><p><br></p><p><strong>chown – Thay đổi chủ sở hữu và nhóm của tập tin.</strong></p><p>chown user:group tên_tập_tin</p><p><br></p><p><strong>umask – Đặt quyền mặc định khi tạo tập tin.</strong></p><p>umask 022&nbsp;</p><p><br></p><h2>Quản Lý Tiến Trình</h2><p><strong>ps – Hiển thị các tiến trình đang chạy.</strong></p><p>ps</p><p>ps aux</p><p><br></p><p><strong>top – Hiển thị thời gian thực các tiến trình và mức sử dụng tài nguyên của hệ thống.</strong></p><p>top</p><p><br></p><p><strong>kill – Dừng tiến trình theo PID của nó.</strong></p><p>kill id_tiến_trình</p><p>kill -9 id_tiến_trình&nbsp;</p><p><br></p><p><strong>htop – Công cụ xem tiến trình tương tác (cần cài đặt).</strong></p><p>htop</p><p><br></p><h2>Thông Tin Hệ Thống</h2><p><strong>df – Hiển thị không gian đĩa sử dụng.</strong></p><p>df -h&nbsp;&nbsp;</p><p><br></p><p><strong>du – Hiển thị dung lượng đĩa sử dụng cho các tập tin và thư mục.</strong></p><p>du -h /đường/dẫn/đến/thư/mục</p><p><br></p><p><strong>free – Hiển thị mức sử dụng bộ nhớ.</strong></p><p>free -h&nbsp;&nbsp;</p><p><br></p><p><strong>uname – Hiển thị thông tin hệ thống.</strong></p><p>uname -a&nbsp;</p><p><br></p><p><strong>uptime – Hiển thị thời gian hệ thống đã hoạt động.</strong></p><p>uptime</p><p><br></p><p><strong>whoami – Hiển thị người dùng hiện đang đăng nhập.</strong></p><p>whoami</p><p><br></p><p><strong>hostname – Hiển thị hoặc đặt tên máy chủ của hệ thống.</strong></p><p>hostname</p><p><br></p><p><strong>lscpu – Hiển thị thông tin kiến trúc CPU.</strong></p><p>lscpu</p><p><br></p><h2>Lệnh Mạng</h2><p><strong>ping – Kiểm tra kết nối với một máy chủ.</strong></p><p>ping google.com</p><p><br></p><p><strong>ifconfig – Hiển thị thông tin giao diện mạng (có thể cần cài đặt net-tools trên một số hệ thống).</strong></p><p>ifconfig</p><p><br></p><p><strong>ip – Cấu hình giao diện mạng và định tuyến.</strong></p><p>ip addr show&nbsp;&nbsp;</p><p>ip route show&nbsp;&nbsp;</p><p><br></p><p><strong>curl – Lấy dữ liệu từ một URL.</strong></p><p>curl https://example.com</p><p><br></p><p><strong>wget – Tải tệp từ web.</strong></p><p>wget https://example.com/file.zip</p><p><br></p><h2>Quản Lý Gói Phần Mềm</h2><p><strong>apt-get (cho các hệ điều hành dựa trên Debian/Ubuntu) – Cài đặt, cập nhật, hoặc gỡ bỏ gói phần mềm.</strong></p><p>sudo apt-get update&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</p><p>sudo apt-get install package&nbsp;</p><p>sudo apt-get remove package&nbsp;</p><p><br></p><p><strong>yum (cho các hệ điều hành dựa trên RedHat/CentOS) – Cài đặt, cập nhật, hoặc gỡ bỏ gói phần mềm.</strong></p><p>sudo yum update&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</p><p>sudo yum install package&nbsp;&nbsp;</p><p>sudo yum remove package&nbsp;&nbsp;</p><p><br></p><p><strong class=\"ql-size-large\" style=\"color: rgb(255, 255, 255);\">Nén &amp; Giải Nén Tập Tin</strong></p><p><strong>tar – Lưu trữ hoặc giải nén các tập tin.</strong></p><p>tar -czvf tên_lưu_trữ.tar.gz /đường/dẫn/đến/thư_mục&nbsp;&nbsp;</p><p>tar -xzvf tên_lưu_trữ.tar.gz&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</p><p><br></p><p><strong>zip – Nén các tập tin thành tệp zip.</strong></p><p>zip tên_lưu_trữ.zip tập_tin1 tập_tin2</p><p><br></p><p><strong>unzip – Giải nén một tệp zip.</strong></p><p>unzip tên_lưu_trữ.zip</p><p><br></p><h2>Lệnh Khác</h2><p><strong>echo – In thông báo hoặc biến ra màn hình terminal.</strong></p><p>echo \"Xin chào, Thế Giới!\"</p><p><br></p><p><strong>date – Hiển thị hoặc đặt ngày và giờ hệ thống.</strong></p><p>date</p><p><br></p><p><strong>alias – Tạo bí danh cho lệnh.</strong></p><p>alias ll=\'ls -la\'&nbsp;&nbsp;</p><p><br></p><p><strong>history – Hiển thị lịch sử lệnh các lệnh vừa được dùng.</strong></p><p>history</p><p><br></p><p><strong>clear – Xóa, làm mới màn hình terminal.</strong></p><p>clear</p>',2,'https://www.extremetech.com/wp-content/uploads/2012/05/Linux-logo-without-version-number-banner-sized.jpg','Linux');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(120) DEFAULT NULL,
  `email` varchar(120) DEFAULT NULL,
  `password_hash` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2025-10-29 09:17:40.518','2025-10-29 09:17:40.518',NULL,'DevOps Maintainer','admin@hocdevops.community','$2a$10$ye/SNFx7K4GFRATVPQBfR.sgLwHFIjLwR.8sWA1kBhBkFjNu0ixxS'),(2,'2025-10-29 09:27:47.416','2025-10-29 09:27:47.416',NULL,'Hiệp','hiep18797@gmail.com','$2a$10$3.0TDcuRlsMf1UL96Y9/eewWjbX1CvuxEnXD9RuXFkYzMHuSOllAK');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-10-31  2:55:06
