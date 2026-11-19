package seed

import (
	"fmt"
	"log"
	"time"

	"nalakarsa/internal/common/constant"
	"nalakarsa/internal/model"
	"nalakarsa/internal/utils"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	log.Println("Seeding database with default mock data...")
	log.Println("Purging existing data...")
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.PublicationSave{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.PublicationCollaborator{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Publication{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.HomepageSection{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.HomepageHero{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.HomepageTestimonial{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ProjectApplication{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.CollaborationRequest{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ProjectMilestone{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ProjectMember{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Project{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Notification{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Message{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ConversationMember{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Conversation{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Connection{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.DiscussionReply{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Discussion{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.RefreshToken{})

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.User{})

	if err := SeedExpertise(db); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return err
	}
	verifiedAt := time.Now().UTC()
	akademisiUser := &model.User{
		Email:           "dosen@nalakarsa.id",
		PasswordHash:    hashedPassword,
		Role:            constant.RoleAkademisi,
		SystemRole:      constant.SystemRoleUser,
		Status:          constant.StatusActive,
		FirstName:       "Budi",
		LastName:        "Santoso",
		FullName:        "Dr. Ir. Budi Santoso",
		PrefixTitle:     "Dr. Ir.",
		SuffixTitle:     "",
		Affiliation:     "Universitas Indonesia",
		Location:        "Jakarta",
		Expertise:       "Kecerdasan Buatan, Machine Learning",
		Bio:             "Peneliti bidang IoT.",
		Mission:         "Melakukan penelitian terapan dan mencari mitra industrialisasi.",
		AvatarURL:       "https://api.dicebear.com/7.x/adventurer/svg?seed=budi",
		EmailVerifiedAt: &verifiedAt,
	}

	praktisiUser := &model.User{
		Email:           "pengusaha@nalakarsa.id",
		PasswordHash:    hashedPassword,
		Role:            constant.RolePraktisi,
		SystemRole:      constant.SystemRoleUser,
		Status:          constant.StatusActive,
		FirstName:       "Setyo",
		LastName:        "Nugroho",
		FullName:        "Setyo Nugroho",
		PrefixTitle:     "",
		SuffixTitle:     "",
		Affiliation:     "PT Telkom Indonesia",
		Location:        "Surabaya",
		Expertise:       "Data Science, Big Data Analytics",
		Bio:             "Pengusaha pertanian digital.",
		Mission:         "Mendukung otomatisasi pertanian menggunakan teknologi hasil riset.",
		AvatarURL:       "https://api.dicebear.com/7.x/adventurer/svg?seed=hendra",
		EmailVerifiedAt: &verifiedAt,
	}

	profesionalUser := &model.User{
		Email:           "engineer@nalakarsa.id",
		PasswordHash:    hashedPassword,
		Role:            constant.RoleProfesional,
		SystemRole:      constant.SystemRoleUser,
		Status:          constant.StatusActive,
		FirstName:       "Hendra",
		LastName:        "Wijaya",
		FullName:        "Hendra Wijaya",
		PrefixTitle:     "",
		SuffixTitle:     "",
		Affiliation:     "Tech Startup Indonesia",
		Location:        "Bandung",
		Expertise:       "Software Engineering, Cloud Architecture",
		Bio:             "Software Engineer spesialis backend.",
		Mission:         "Menyediakan solusi backend handal skala besar.",
		AvatarURL:       "https://api.dicebear.com/7.x/adventurer/svg?seed=setyo",
		EmailVerifiedAt: &verifiedAt,
	}

	if err := db.Create(akademisiUser).Error; err != nil {
		return err
	}
	if err := db.Create(praktisiUser).Error; err != nil {
		return err
	}
	if err := db.Create(profesionalUser).Error; err != nil {
		return err
	}

	discussionUsers := []*model.User{
		{
			Email: "peneliti01@nalakarsa.id", PasswordHash: hashedPassword, Role: constant.RoleAkademisi, SystemRole: constant.SystemRoleUser, Status: constant.StatusActive,
			FirstName: "Rina", LastName: "Kusuma", FullName: "Rina Kusuma", Affiliation: "Institut Teknologi Bandung", Location: "Bandung", Expertise: "Energi Terbarukan, IoT", Bio: "Peneliti energi terbarukan dan sistem monitoring.", Mission: "Mendorong pemanfaatan teknologi untuk transisi energi.", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=rina", EmailVerifiedAt: &verifiedAt,
		},
		{
			Email: "peneliti02@nalakarsa.id", PasswordHash: hashedPassword, Role: constant.RoleAkademisi, SystemRole: constant.SystemRoleUser, Status: constant.StatusActive,
			FirstName: "Dimas", LastName: "Pratama", FullName: "Dimas Pratama", Affiliation: "Universitas Gadjah Mada", Location: "Yogyakarta", Expertise: "Pendidikan Digital, NLP", Bio: "Dosen yang mengembangkan pembelajaran berbasis data.", Mission: "Membuat pendidikan lebih adaptif dan inklusif.", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=dimas", EmailVerifiedAt: &verifiedAt,
		},
		{
			Email: "praktisi01@nalakarsa.id", PasswordHash: hashedPassword, Role: constant.RolePraktisi, SystemRole: constant.SystemRoleUser, Status: constant.StatusActive,
			FirstName: "Lestari", LastName: "Wulandari", FullName: "Lestari Wulandari", Affiliation: "AgriTech Nusantara", Location: "Bogor", Expertise: "Agrikultur, Supply Chain", Bio: "Praktisi agribisnis yang fokus pada efisiensi rantai pasok.", Mission: "Menghubungkan inovasi pertanian dengan kebutuhan petani.", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=lestari", EmailVerifiedAt: &verifiedAt,
		},
		{
			Email: "praktisi02@nalakarsa.id", PasswordHash: hashedPassword, Role: constant.RolePraktisi, SystemRole: constant.SystemRoleUser, Status: constant.StatusActive,
			FirstName: "Bagas", LastName: "Ramadhan", FullName: "Bagas Ramadhan", Affiliation: "Manufaktur Maju", Location: "Bekasi", Expertise: "Manufaktur, Otomasi", Bio: "Manajer inovasi di sektor manufaktur.", Mission: "Meningkatkan produktivitas industri melalui otomasi.", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=bagas", EmailVerifiedAt: &verifiedAt,
		},
		{
			Email: "profesional01@nalakarsa.id", PasswordHash: hashedPassword, Role: constant.RoleProfesional, SystemRole: constant.SystemRoleUser, Status: constant.StatusActive,
			FirstName: "Nadia", LastName: "Permata", FullName: "Nadia Permata", Affiliation: "Data Insight Lab", Location: "Jakarta", Expertise: "Data Engineering, Analytics", Bio: "Data engineer yang membantu organisasi mengambil keputusan berbasis data.", Mission: "Membangun fondasi data yang mudah digunakan dan terpercaya.", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=nadia", EmailVerifiedAt: &verifiedAt,
		},
		{
			Email: "profesional02@nalakarsa.id", PasswordHash: hashedPassword, Role: constant.RoleProfesional, SystemRole: constant.SystemRoleUser, Status: constant.StatusActive,
			FirstName: "Yoga", LastName: "Mahendra", FullName: "Yoga Mahendra", Affiliation: "Cloudworks Indonesia", Location: "Surabaya", Expertise: "Cloud Computing, DevOps", Bio: "Cloud architect untuk platform yang aman dan skalabel.", Mission: "Membantu tim membangun layanan digital yang tangguh.", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=yoga", EmailVerifiedAt: &verifiedAt,
		},
		{
			Email: "profesional03@nalakarsa.id", PasswordHash: hashedPassword, Role: constant.RoleProfesional, SystemRole: constant.SystemRoleUser, Status: constant.StatusActive,
			FirstName: "Salsa", LastName: "Anindita", FullName: "Salsa Anindita", Affiliation: "HealthTech Indonesia", Location: "Semarang", Expertise: "Teknologi Kesehatan, Product Design", Bio: "Product specialist untuk solusi kesehatan digital.", Mission: "Membuat layanan kesehatan lebih mudah dijangkau.", AvatarURL: "https://api.dicebear.com/7.x/adventurer/svg?seed=salsa", EmailVerifiedAt: &verifiedAt,
		},
	}
	for _, user := range discussionUsers {
		if err := db.Create(user).Error; err != nil {
			return err
		}
	}
	authors := append([]*model.User{akademisiUser, praktisiUser, profesionalUser}, discussionUsers...)
	if err := seedPublications(db, akademisiUser, praktisiUser, profesionalUser); err != nil {
		return err
	}
	discussions, err := seedDiscussions(db, authors)
	if err != nil {
		return err
	}
	disc1, disc2 := discussions[0], discussions[1]
	reply1 := &model.DiscussionReply{
		DiscussionID: disc1.ID,
		UserID:       akademisiUser.ID,
		Content:      "Setuju sekali! Kami di lab riset juga menggunakan Go untuk backend sensor gateway kami karena penggunaan memori yang sangat kecil.",
	}

	reply2 := &model.DiscussionReply{
		DiscussionID: disc2.ID,
		UserID:       praktisiUser.ID,
		Content:      "Kami dari sektor swasta siap menampung riset yang sudah matang di TRL 7 ke atas untuk dikomersialkan. Silakan hubungi kami.",
	}

	if err := db.Create(reply1).Error; err != nil {
		return err
	}
	if err := db.Create(reply2).Error; err != nil {
		return err
	}
	proj := &model.Project{
		OwnerID:       akademisiUser.ID,
		Title:         "Pengembangan Lapangan Sistem IoT Monitoring Tanah Pertanian",
		Description:   "Kami merancang sensor kelembapan tanah berbasis LoRaWAN. Kami mencari Praktisi Agribisnis yang memiliki lahan perkebunan untuk melakukan uji coba lapangan nyata dan validasi pasar.",
		Category:      "Research",
		Needs:         constant.RolePraktisi,
		Status:        constant.ProjectStatusActive,
		Progress:      45,
		FundingStatus: "Self-funded",
		Location:      "Bogor",
		Deadline:      timePtr(time.Now().AddDate(0, 2, 0)),
	}

	if err := db.Create(proj).Error; err != nil {
		return err
	}
	milestones := []model.ProjectMilestone{
		{
			ProjectID:  proj.ID,
			Title:      "Validasi kebutuhan lahan uji",
			DueAt:      timePtr(time.Now().AddDate(0, 0, 14)),
			Status:     constant.MilestoneStatusCompleted,
			AssigneeID: &praktisiUser.ID,
		},
		{
			ProjectID:  proj.ID,
			Title:      "Pemasangan sensor dan pengumpulan data awal",
			DueAt:      timePtr(time.Now().AddDate(0, 0, 30)),
			Status:     constant.MilestoneStatusInProgress,
			AssigneeID: &praktisiUser.ID,
		},
	}
	if err := db.Create(&milestones).Error; err != nil {
		return err
	}
	if err := db.Create(&model.ProjectMember{
		ProjectID: proj.ID,
		UserID:    praktisiUser.ID,
		Role:      "Mitra Agribisnis",
		Status:    constant.MemberStatusActive,
	}).Error; err != nil {
		return err
	}
	app := &model.ProjectApplication{
		ProjectID:   proj.ID,
		ApplicantID: praktisiUser.ID,
		Message:     "Halo Pak Budi, saya Hendra dari PT Tani Maju Digital. Kami memiliki perkebunan teh seluas 3 hektar di daerah Bogor dan sangat tertarik untuk dijadikan lahan riset implementasi alat ini.",
		Status:      constant.ApplicationStatusPending,
	}

	if err := db.Create(app).Error; err != nil {
		return err
	}
	if err := seedCollaborationRequests(db, akademisiUser, praktisiUser, profesionalUser, proj, disc1, disc2); err != nil {
		return err
	}
	if err := db.Create(&model.Connection{
		RequesterID: akademisiUser.ID,
		AddresseeID: praktisiUser.ID,
		Status:      constant.ConnectionStatusAccepted,
	}).Error; err != nil {
		return err
	}
	now := time.Now()
	conversation := &model.Conversation{Type: constant.ConversationTypeDirect, LastMessageAt: &now}
	if err := db.Create(conversation).Error; err != nil {
		return err
	}
	if err := db.Create(&[]model.ConversationMember{
		{ConversationID: conversation.ID, UserID: akademisiUser.ID},
		{ConversationID: conversation.ID, UserID: praktisiUser.ID},
	}).Error; err != nil {
		return err
	}
	if err := db.Create(&model.Message{
		ConversationID: conversation.ID,
		SenderID:       praktisiUser.ID,
		Body:           "Halo Pak Budi, saya siap membantu validasi lapangan untuk proyek ini.",
		Status:         "sent",
		CreatedAt:      now,
		UpdatedAt:      now,
	}).Error; err != nil {
		return err
	}
	if err := db.Create(&model.Notification{
		UserID:       akademisiUser.ID,
		ActorID:      &praktisiUser.ID,
		Type:         constant.TypeConnectionAccepted,
		ResourceType: constant.ResourceConnection,
		Payload:      "{}",
		CreatedAt:    now,
	}).Error; err != nil {
		return err
	}
	hero := &model.HomepageHero{
		Headline:     "Kolaborasi Riset jadi Solusi Nyata",
		SubHeadline:  "Temukan akademisi, praktisi, dan profesional untuk mendorong riset yang berdampak.",
		CallToAction: "Jelajahi Diskusi",
		CTAURL:       "/discussions",
		ImageURL:     "/assets/homepage/hero.jpg",
		IsActive:     true,
	}

	sections := []model.HomepageSection{
		{
			Key:       "benefits",
			Title:     "Mengapa Nalakarsa?",
			Subtitle:  "Platform penghubung ekosistem riset dan industri.",
			Content:   "Temukan kolaborasi lintas peran untuk mengubah ide riset menjadi implementasi nyata.",
			ImageURL:  "/assets/homepage/benefits.jpg",
			LinkLabel: "Pelajari Selengkapnya",
			LinkURL:   "/about",
			SortOrder: 1,
			IsActive:  true,
		},
		{
			Key:       "highlights",
			Title:     "Proyek & Kolaborasi",
			Subtitle:  "Akses peluang yang relevan.",
			Content:   "Lihat proyek, ikuti diskusi, dan bangun koneksi yang membantu kemajuan industri.",
			ImageURL:  "/assets/homepage/highlights.jpg",
			LinkLabel: "Lihat Proyek",
			LinkURL:   "/projects",
			SortOrder: 2,
			IsActive:  true,
		},
		{
			Key:       "research-network",
			Title:     "Bangun Jaringan Riset",
			Subtitle:  "Terhubung dengan orang yang tepat.",
			Content:   "Temukan akademisi, praktisi, dan profesional berdasarkan keahlian serta kebutuhan kolaborasi.",
			ImageURL:  "/assets/homepage/research-network.jpg",
			LinkLabel: "Cari Kolaborator",
			LinkURL:   "/users",
			SortOrder: 3,
			IsActive:  true,
		},
		{
			Key:       "publication-feed",
			Title:     "Bagikan Pengetahuan",
			Subtitle:  "Publikasikan karya dan hasil riset.",
			Content:   "Jadikan hasil penelitian lebih mudah ditemukan dan dikembangkan bersama oleh ekosistem inovasi.",
			ImageURL:  "/assets/homepage/publication-feed.jpg",
			LinkLabel: "Jelajahi Publikasi",
			LinkURL:   "/publications",
			SortOrder: 4,
			IsActive:  true,
		},
		{
			Key:       "impact",
			Title:     "Dorong Dampak Nyata",
			Subtitle:  "Dari ide menjadi solusi.",
			Content:   "Kelola kolaborasi, diskusi, dan proyek dalam satu platform untuk mempercepat hilirisasi inovasi.",
			ImageURL:  "/assets/homepage/impact.jpg",
			LinkLabel: "Mulai Berkolaborasi",
			LinkURL:   "/auth/register",
			SortOrder: 5,
			IsActive:  true,
		},
	}

	testimonials := []model.HomepageTestimonial{
		{
			Name:      "Dr. Raka Santoso",
			Role:      "Dosen",
			Company:   "Universitas Teknologi Nusantara",
			Message:   "Nalakarsa membantu saya menemukan mitra untuk mengimplementasikan riset ke pilot project.",
			AvatarURL: "/assets/testimonials/raka.jpg",
			IsActive:  true,
			SortOrder: 1,
		},
		{
			Name:      "Sari Wijaya",
			Role:      "Praktisi",
			Company:   "AgriTech Indonesia",
			Message:   "Kolaborasi jadi lebih cepat karena diskusi dan proyek terstruktur dalam satu platform.",
			AvatarURL: "/assets/testimonials/sari.jpg",
			IsActive:  true,
			SortOrder: 2,
		},
		{
			Name:      "Andi Pratama",
			Role:      "Profesional",
			Company:   "Teknologi Nusantara",
			Message:   "Saya dapat menemukan proyek yang sesuai dengan keahlian dan ikut berkontribusi sejak tahap perencanaan.",
			AvatarURL: "/assets/testimonials/andi.jpg",
			IsActive:  true,
			SortOrder: 3,
		},
		{
			Name:      "Maya Lestari",
			Role:      "Peneliti",
			Company:   "Pusat Inovasi Indonesia",
			Message:   "Nalakarsa membantu hasil riset kami menjangkau mitra industri yang tepat dan siap berkolaborasi.",
			AvatarURL: "/assets/testimonials/maya.jpg",
			IsActive:  true,
			SortOrder: 4,
		},
	}

	if err := db.Create(hero).Error; err != nil {
		return err
	}
	if err := db.Create(&sections).Error; err != nil {
		return err
	}
	if err := db.Create(&testimonials).Error; err != nil {
		return err
	}

	log.Println("Mock data seeding successfully completed.")
	return nil
}

func seedDiscussions(db *gorm.DB, authors []*model.User) ([]*model.Discussion, error) {
	definitions := []struct {
		title, description, background, objective, expectedOutput, category, tags string
	}{
		{"Mengapa Golang Sangat Cocok untuk Microservices?", "Diskusi tentang konkurensi, goroutines, dan channels untuk layanan backend yang efisien.", "Tim membutuhkan arsitektur layanan yang ringan dan mudah dikembangkan.", "Membandingkan praktik terbaik microservices berbasis Go.", "Rancangan arsitektur microservices", "Teknologi", "golang,microservices,backend"},
		{"Hilirisasi Hasil Riset Universitas ke Sektor Industri", "Bagaimana menjembatani prototipe riset agar dapat digunakan oleh pelaku industri?", "Banyak prototipe riset berhenti di laboratorium karena kurangnya mitra implementasi.", "Menemukan pola kolaborasi antara kampus dan industri.", "Model kerja sama hilirisasi", "Riset", "riset,industri,hilirisasi"},
		{"Monitoring Kualitas Tanah Berbasis IoT", "Membahas sensor dan dashboard untuk membantu petani mengambil keputusan berbasis data.", "Petani membutuhkan informasi kondisi tanah yang cepat dan mudah dipahami.", "Merancang pilot project monitoring tanah.", "Proposal uji lapangan sensor", "Pertanian", "iot,pertanian,sensor"},
		{"Platform Pembelajaran Adaptif untuk Sekolah", "Berbagi pengalaman membangun sistem pembelajaran yang menyesuaikan kebutuhan siswa.", "Materi dan metode belajar perlu menyesuaikan kemampuan setiap siswa.", "Mengidentifikasi data yang aman dan bermanfaat untuk personalisasi belajar.", "Konsep platform pembelajaran adaptif", "Pendidikan", "pendidikan,edutech,data"},
		{"Interoperabilitas Data Kesehatan Nasional", "Mencari pendekatan untuk menghubungkan data fasilitas kesehatan tanpa mengabaikan privasi pasien.", "Data kesehatan masih tersebar di banyak sistem dengan format berbeda.", "Menyusun prinsip integrasi data kesehatan yang aman.", "Prinsip arsitektur interoperabilitas", "Kesehatan", "kesehatan,interoperabilitas,privasi"},
		{"Pemanfaatan Energi Surya untuk UMKM", "Mendiskusikan model adopsi panel surya yang terjangkau bagi usaha kecil.", "Biaya energi menjadi salah satu beban operasional UMKM.", "Mengukur kelayakan solusi energi surya untuk UMKM.", "Studi kelayakan awal", "Energi", "energi-surya,umkm,keberlanjutan"},
		{"Otomasi Proses Produksi yang Terjangkau", "Berbagi strategi memulai otomasi pada pabrik skala menengah.", "Pelaku manufaktur membutuhkan otomasi bertahap dengan investasi terukur.", "Menyusun prioritas proses yang layak diotomasi.", "Roadmap otomasi produksi", "Manufaktur", "manufaktur,otomasi,produktivitas"},
		{"Analitik Data untuk Rantai Pasok Pangan", "Membahas penggunaan analitik untuk mengurangi pemborosan dan keterlambatan distribusi pangan.", "Perubahan permintaan dan distribusi membuat rantai pasok sulit diprediksi.", "Menentukan indikator analitik yang paling berguna bagi pelaku usaha.", "Daftar indikator supply chain", "Ekonomi", "data-analytics,supply-chain,pangan"},
		{"Kebijakan Data Terbuka untuk Inovasi Daerah", "Bagaimana pemerintah daerah dapat membuka data tanpa mengorbankan keamanan dan privasi?", "Data publik yang mudah digunakan dapat mempercepat inovasi lokal.", "Menyusun rekomendasi tata kelola data terbuka.", "Draft rekomendasi kebijakan", "Kebijakan Publik", "data-terbuka,kebijakan,inovasi"},
		{"Kolaborasi Teknologi untuk Kota Berkelanjutan", "Mencari ide kolaborasi lintas bidang untuk mengurangi dampak lingkungan perkotaan.", "Pertumbuhan kota meningkatkan kebutuhan solusi transportasi dan pengelolaan lingkungan.", "Mengumpulkan ide proyek kota berkelanjutan yang dapat diuji.", "Daftar ide pilot project", "Lingkungan", "smart-city,lingkungan,kolaborasi"},
	}

	if len(authors) < len(definitions) {
		return nil, fmt.Errorf("discussion seeder requires %d authors, got %d", len(definitions), len(authors))
	}

	discussions := make([]*model.Discussion, 0, len(definitions))
	for i, definition := range definitions {
		discussion := &model.Discussion{
			UserID: authors[i].ID, Title: definition.title, Description: definition.description,
			Background: definition.background, Objective: definition.objective, ExpectedOutput: definition.expectedOutput,
			Category: definition.category, Tags: definition.tags, Status: "open",
		}
		if err := db.Create(discussion).Error; err != nil {
			return nil, err
		}
		discussions = append(discussions, discussion)
	}
	return discussions, nil
}

// seedCollaborationRequests creates representative records for the
// collaboration-request endpoints, including requests from discussions and
// projects in each supported status.
func seedCollaborationRequests(db *gorm.DB, akademisi, praktisi, profesional *model.User, project *model.Project, disc1, disc2 *model.Discussion) error {
	rejectionReason := "Kontribusi belum sesuai dengan kebutuhan proyek saat ini."
	requests := []model.CollaborationRequest{
		{DiscussionID: &disc1.ID, ApplicantID: akademisi.ID, ProposedContribution: "Menyusun rancangan arsitektur microservices dan dokumentasi teknis.", Status: constant.CollabStatusPending},
		{DiscussionID: &disc1.ID, ApplicantID: praktisi.ID, ProposedContribution: "Menyediakan studi kasus implementasi backend di lingkungan produksi.", Status: constant.CollabStatusPending},
		{DiscussionID: &disc1.ID, ApplicantID: profesional.ID, ProposedContribution: "Membantu review performa, observability, dan deployment layanan.", Status: constant.CollabStatusAccepted},
		{DiscussionID: &disc2.ID, ApplicantID: praktisi.ID, ProposedContribution: "Membuka akses pilot project untuk validasi teknologi di lapangan.", Status: constant.CollabStatusPending},
		{DiscussionID: &disc2.ID, ApplicantID: profesional.ID, ProposedContribution: "Membuat dashboard monitoring hasil hilirisasi riset.", Status: constant.CollabStatusRejected, RejectionReason: &rejectionReason},
		{ProjectID: &project.ID, ApplicantID: akademisi.ID, ProposedContribution: "Menyediakan konsultasi metodologi penelitian dan validasi data.", Status: constant.CollabStatusPending},
		{ProjectID: &project.ID, ApplicantID: praktisi.ID, ProposedContribution: "Menyiapkan lokasi uji dan pendampingan operasional lapangan.", Status: constant.CollabStatusAccepted},
		{ProjectID: &project.ID, ApplicantID: profesional.ID, ProposedContribution: "Mengembangkan API integrasi sensor dan pipeline data.", Status: constant.CollabStatusPending},
		{ProjectID: &project.ID, ApplicantID: akademisi.ID, ProposedContribution: "Membantu analisis hasil pengujian dan penyusunan laporan riset.", Status: constant.CollabStatusPending},
		{ProjectID: &project.ID, ApplicantID: profesional.ID, ProposedContribution: "Menyusun rencana keamanan aplikasi dan deployment cloud.", Status: constant.CollabStatusPending},
	}

	return db.Create(&requests).Error
}

func seedPublications(db *gorm.DB, users ...*model.User) error {
	if len(users) < 3 {
		return nil
	}
	types := []string{"CASE_STUDY", "PRACTICAL_ARTICLE", "APPLIED_RESEARCH", "POLICY_BRIEF", "IMPLEMENTATION_REPORT", "IMPLEMENTATION_GUIDE"}
	categories := []string{"Teknologi", "Pendidikan", "Pertanian", "Kebijakan Publik", "Kesehatan", "Industri"}
	titles := []string{
		"Implementasi IoT untuk Monitoring Tanah Pertanian",
		"Panduan Praktis Hilirisasi Hasil Riset",
		"Model Kolaborasi Riset Terapan Perguruan Tinggi",
		"Kebijakan Data Terbuka untuk Inovasi Daerah",
		"Laporan Implementasi Sistem Informasi Kesehatan",
		"Panduan Implementasi Transformasi Digital UMKM",
	}
	now := time.Now().UTC()
	for i := range titles {
		publication := model.Publication{
			OwnerID: users[i%len(users)].ID, PublicationType: types[i], Title: titles[i],
			Summary:  "Publikasi contoh untuk menguji publication feed dan kolaborasi Nalakarsa.",
			Content:  "Konten publikasi contoh. Ganti dengan naskah dan hasil implementasi yang sebenarnya.",
			Category: categories[i], Status: "PUBLISHED", CoverURL: "", PublishedAt: timePtr(now.Add(-time.Duration(i) * 24 * time.Hour)),
		}
		if err := db.Create(&publication).Error; err != nil {
			return err
		}
		for j, user := range users {
			if err := db.Create(&model.PublicationCollaborator{PublicationID: publication.ID, UserID: user.ID, ContributionRole: []string{"Penulis", "Mitra", "Reviewer"}[j], IsOwner: user.ID == publication.OwnerID}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func timePtr(value time.Time) *time.Time {
	return &value
}
