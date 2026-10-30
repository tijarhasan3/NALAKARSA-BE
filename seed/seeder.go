package seed

import (
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
	akademisiUser := &model.User{
		Email:        "dosen@nalakarsa.id",
		PasswordHash: hashedPassword,
		Role:         constant.RoleAkademisi,
		SystemRole:   constant.SystemRoleUser,
		Status:       constant.StatusActive,
		FirstName:    "Budi",
		LastName:     "Santoso",
		FullName:     "Dr. Ir. Budi Santoso",
		PrefixTitle:  "Dr. Ir.",
		SuffixTitle:  "",
		Affiliation:  "Universitas Indonesia",
		Location:     "Jakarta",
		Expertise:    "Kecerdasan Buatan, Machine Learning",
		Bio:          "Peneliti bidang IoT.",
		Mission:      "Melakukan penelitian terapan dan mencari mitra industrialisasi.",
		AvatarURL:    "https://api.dicebear.com/7.x/adventurer/svg?seed=budi",
	}

	praktisiUser := &model.User{
		Email:        "pengusaha@nalakarsa.id",
		PasswordHash: hashedPassword,
		Role:         constant.RolePraktisi,
		SystemRole:   constant.SystemRoleUser,
		Status:       constant.StatusActive,
		FirstName:    "Setyo",
		LastName:     "Nugroho",
		FullName:     "Setyo Nugroho",
		PrefixTitle:  "",
		SuffixTitle:  "",
		Affiliation:  "PT Telkom Indonesia",
		Location:     "Surabaya",
		Expertise:    "Data Science, Big Data Analytics",
		Bio:          "Pengusaha pertanian digital.",
		Mission:      "Mendukung otomatisasi pertanian menggunakan teknologi hasil riset.",
		AvatarURL:    "https://api.dicebear.com/7.x/adventurer/svg?seed=hendra",
	}

	profesionalUser := &model.User{
		Email:        "engineer@nalakarsa.id",
		PasswordHash: hashedPassword,
		Role:         constant.RoleProfesional,
		SystemRole:   constant.SystemRoleUser,
		Status:       constant.StatusActive,
		FirstName:    "Hendra",
		LastName:     "Wijaya",
		FullName:     "Hendra Wijaya",
		PrefixTitle:  "",
		SuffixTitle:  "",
		Affiliation:  "Tech Startup Indonesia",
		Location:     "Bandung",
		Expertise:    "Software Engineering, Cloud Architecture",
		Bio:          "Software Engineer spesialis backend.",
		Mission:      "Menyediakan solusi backend handal skala besar.",
		AvatarURL:    "https://api.dicebear.com/7.x/adventurer/svg?seed=setyo",
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
	if err := seedPublications(db, akademisiUser, praktisiUser, profesionalUser); err != nil {
		return err
	}
	disc1 := &model.Discussion{
		UserID:      profesionalUser.ID,
		Title:       "Mengapa Golang Sangat Cocok untuk Microservices?",
		Description: "Dalam ekosistem Nalakarsa, saya ingin berdiskusi mengenai konkurensi di Go. Goroutines dan Channels membuat pemrosesan paralel menjadi sangat efisien dibandingkan thread konvensional. Bagaimana pengalaman rekan-rekan?",
		Category:    "Tech",
		Tags:        "golang,microservices,backend",
		Status:      "open",
	}

	disc2 := &model.Discussion{
		UserID:      akademisiUser.ID,
		Title:       "Hilirisasi Hasil Riset Universitas ke Sektor Industri",
		Description: "Banyak prototipe riset IoT terhenti di laci laboratorium. Kendala utama adalah kurangnya komunikasi dengan pelaku bisnis. Mari kita bahas bagaimana menjembatani gap ini.",
		Category:    "Research",
		Tags:        "riset,industri,hilirisasi",
		Status:      "open",
	}

	if err := db.Create(disc1).Error; err != nil {
		return err
	}
	if err := db.Create(disc2).Error; err != nil {
		return err
	}
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
