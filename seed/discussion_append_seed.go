package seed

import (
	"fmt"
	"time"

	"nalakarsa/internal/common/constant"
	"nalakarsa/internal/model"
	"nalakarsa/internal/utils"

	"gorm.io/gorm"
)

// SeedDiscussionData appends demo discussion data without deleting or
// overwriting existing users and discussions.
func SeedDiscussionData(db *gorm.DB) error {
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return err
	}
	verifiedAt := time.Now().UTC()

	users := []model.User{
		{Email: "dosen@nalakarsa.id", Role: constant.RoleAkademisi, FirstName: "Budi", LastName: "Santoso", FullName: "Dr. Ir. Budi Santoso", Affiliation: "Universitas Indonesia", Location: "Jakarta", Expertise: "Kecerdasan Buatan, Machine Learning", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=budi"},
		{Email: "pengusaha@nalakarsa.id", Role: constant.RolePraktisi, FirstName: "Setyo", LastName: "Nugroho", FullName: "Setyo Nugroho", Affiliation: "PT Telkom Indonesia", Location: "Surabaya", Expertise: "Data Science, Big Data Analytics", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=setyo"},
		{Email: "engineer@nalakarsa.id", Role: constant.RoleProfesional, FirstName: "Hendra", LastName: "Wijaya", FullName: "Hendra Wijaya", Affiliation: "Tech Startup Indonesia", Location: "Bandung", Expertise: "Software Engineering, Cloud Architecture", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=hendra"},
		{Email: "peneliti01@nalakarsa.id", Role: constant.RoleAkademisi, FirstName: "Rina", LastName: "Kusuma", FullName: "Rina Kusuma", Affiliation: "Institut Teknologi Bandung", Location: "Bandung", Expertise: "Energi Terbarukan, IoT", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=rina"},
		{Email: "peneliti02@nalakarsa.id", Role: constant.RoleAkademisi, FirstName: "Dimas", LastName: "Pratama", FullName: "Dimas Pratama", Affiliation: "Universitas Gadjah Mada", Location: "Yogyakarta", Expertise: "Pendidikan Digital, NLP", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=dimas"},
		{Email: "praktisi01@nalakarsa.id", Role: constant.RolePraktisi, FirstName: "Lestari", LastName: "Wulandari", FullName: "Lestari Wulandari", Affiliation: "AgriTech Nusantara", Location: "Bogor", Expertise: "Agrikultur, Supply Chain", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=lestari"},
		{Email: "praktisi02@nalakarsa.id", Role: constant.RolePraktisi, FirstName: "Bagas", LastName: "Ramadhan", FullName: "Bagas Ramadhan", Affiliation: "Manufaktur Maju", Location: "Bekasi", Expertise: "Manufaktur, Otomasi", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=bagas"},
		{Email: "profesional01@nalakarsa.id", Role: constant.RoleProfesional, FirstName: "Nadia", LastName: "Permata", FullName: "Nadia Permata", Affiliation: "Data Insight Lab", Location: "Jakarta", Expertise: "Data Engineering, Analytics", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=nadia"},
		{Email: "profesional02@nalakarsa.id", Role: constant.RoleProfesional, FirstName: "Yoga", LastName: "Mahendra", FullName: "Yoga Mahendra", Affiliation: "Cloudworks Indonesia", Location: "Surabaya", Expertise: "Cloud Computing, DevOps", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=yoga"},
		{Email: "profesional03@nalakarsa.id", Role: constant.RoleProfesional, FirstName: "Salsa", LastName: "Anindita", FullName: "Salsa Anindita", Affiliation: "HealthTech Indonesia", Location: "Semarang", Expertise: "Teknologi Kesehatan, Product Design", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=salsa"},
	}

	authors := make([]*model.User, 0, len(users))
	for i := range users {
		users[i].PasswordHash = hashedPassword
		users[i].SystemRole = constant.SystemRoleUser
		users[i].Status = constant.StatusActive
		users[i].EmailVerifiedAt = &verifiedAt
		if err := db.Where("email = ?", users[i].Email).FirstOrCreate(&users[i]).Error; err != nil {
			return fmt.Errorf("seed discussion user %s: %w", users[i].Email, err)
		}
		authors = append(authors, &users[i])
	}

	discussions := []struct {
		title, description, category, tags string
	}{
		{"Mengapa Golang Sangat Cocok untuk Microservices?", "Diskusi tentang konkurensi, goroutines, dan channels untuk layanan backend yang efisien.", "Teknologi", "golang,microservices,backend"},
		{"Hilirisasi Hasil Riset Universitas ke Sektor Industri", "Bagaimana menjembatani prototipe riset agar dapat digunakan oleh pelaku industri?", "Riset", "riset,industri,hilirisasi"},
		{"Monitoring Kualitas Tanah Berbasis IoT", "Membahas sensor dan dashboard untuk membantu petani mengambil keputusan berbasis data.", "Pertanian", "iot,pertanian,sensor"},
		{"Platform Pembelajaran Adaptif untuk Sekolah", "Berbagi pengalaman membangun sistem pembelajaran yang menyesuaikan kebutuhan siswa.", "Pendidikan", "pendidikan,edutech,data"},
		{"Interoperabilitas Data Kesehatan Nasional", "Mencari pendekatan untuk menghubungkan data fasilitas kesehatan tanpa mengabaikan privasi pasien.", "Kesehatan", "kesehatan,interoperabilitas,privasi"},
		{"Pemanfaatan Energi Surya untuk UMKM", "Mendiskusikan model adopsi panel surya yang terjangkau bagi usaha kecil.", "Energi", "energi-surya,umkm,keberlanjutan"},
		{"Otomasi Proses Produksi yang Terjangkau", "Berbagi strategi memulai otomasi pada pabrik skala menengah.", "Manufaktur", "manufaktur,otomasi,produktivitas"},
		{"Analitik Data untuk Rantai Pasok Pangan", "Membahas penggunaan analitik untuk mengurangi pemborosan distribusi pangan.", "Ekonomi", "data-analytics,supply-chain,pangan"},
		{"Kebijakan Data Terbuka untuk Inovasi Daerah", "Bagaimana pemerintah daerah dapat membuka data tanpa mengorbankan keamanan dan privasi?", "Kebijakan Publik", "data-terbuka,kebijakan,inovasi"},
		{"Kolaborasi Teknologi untuk Kota Berkelanjutan", "Mencari ide kolaborasi lintas bidang untuk mengurangi dampak lingkungan perkotaan.", "Lingkungan", "smart-city,lingkungan,kolaborasi"},
	}

	for i, definition := range discussions {
		discussion := &model.Discussion{
			UserID: authors[i].ID, Title: definition.title, Description: definition.description,
			Category: definition.category, Tags: definition.tags, Status: "open",
		}
		if err := db.Where("user_id = ? AND title = ?", discussion.UserID, discussion.Title).FirstOrCreate(discussion).Error; err != nil {
			return fmt.Errorf("seed discussion %q: %w", discussion.Title, err)
		}
	}
	return nil
}
