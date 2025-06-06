package seed

import (
	"content-management-service/data_layer/entity"
	"log"

	"gorm.io/gorm"
)

// ReaderSeeder seeds the readers table
func ReaderSeeder(db *gorm.DB) error {
	log.Println("Seeding readers...")

	readers := []entity.Reader{
		{Name: "Karen Savage"},
		{Name: "LibriVox Volunteers"},
		{Name: "Group"},
		{Name: "Ancilla"},
		{Name: "Elizabeth Klett"},
		{Name: "Andrea Fiore"},
		{Name: "Phil Chenevert"},
		{Name: "Eden Rea-Hedrick"},
		{Name: "Bob Neufeld"},
		{Name: "John W. Michaels"},
		{Name: "Meredith Hughes"},
		{Name: "Sue Anderson"},
		{Name: "Ruth Golding"},
		{Name: "KirksVoice"},
		{Name: "Hugh McGuire"},
		{Name: "Laurie Anne Walden"},
		{Name: "Tadhg"},
		{Name: "Cori Samuel"},
		{Name: "Mark F. Smith"},
		{Name: "Lars Rolander (1942-2016)"},
		{Name: "Lee Smalley"},
		{Name: "Nan Dodge"},
		{Name: "Ted Delorme"},
		{Name: "David Richardson"},
		{Name: "rachelellen"},
		{Name: "Rob De Lorenzo"},
		{Name: "Adrian Praetzellis"},
		{Name: "DavidG"},
		{Name: "Steve C"},
		{Name: "Peter John Keeble"},
		{Name: "David Wales"},
		{Name: "Kate Follis"},
		{Name: "Simon Evers"},
		{Name: "Phil Benson"},
		{Name: "SweetHome"},
		{Name: "Jim Locke"},
		{Name: "Scott Miller"},
		{Name: "Sarah Mitchell"},
		{Name: "Rachel Thomson"},
		{Name: "Michael Chen"},
		{Name: "Jennifer Walsh"},
		{Name: "Robert Garcia"},
		{Name: "Lisa Parker"},
		{Name: "Thomas Anderson"},
		{Name: "Maria Rodriguez"},
		{Name: "Daniel Kim"},
		{Name: "Susan Brown"},
		{Name: "James Wilson"},
		{Name: "Amanda Taylor"},
		{Name: "Christopher Lee"},
		{Name: "Michelle Johnson"},
		{Name: "Kevin Davis"},
		{Name: "Laura Martinez"},
		{Name: "Brian Thompson"},
		{Name: "Nicole White"},
		{Name: "Steven Clark"},
		{Name: "Jessica Lewis"},
		{Name: "Mark Williams"},
		{Name: "Stephanie Jones"},
		{Name: "Matthew Miller"},
		{Name: "Ashley Garcia"},
		{Name: "Andrew Robinson"},
		{Name: "Samantha Hall"},
		{Name: "Joshua Young"},
		{Name: "Emily Wright"},
		{Name: "Ryan King"},
		{Name: "Melissa Scott"},
		{Name: "Nicholas Green"},
		{Name: "Rebecca Baker"},
		{Name: "Jonathan Adams"},
		{Name: "Katherine Nelson"},
		{Name: "Adam Carter"},
		{Name: "Brittany Mitchell"},
		{Name: "Brandon Roberts"},
		{Name: "Danielle Turner"},
		{Name: "Tyler Phillips"},
		{Name: "Vanessa Campbell"},
		{Name: "Kyle Parker"},
		{Name: "Christina Evans"},
		{Name: "Jeremy Edwards"},
		{Name: "Natalie Collins"},
		{Name: "Austin Stewart"},
		{Name: "Monica Sanchez"},
		{Name: "Sean Morris"},
		{Name: "Tiffany Rogers"},
		{Name: "Derek Reed"},
		{Name: "Heather Cook"},
		{Name: "Ian Bailey"},
		{Name: "Kimberly Rivera"},
		{Name: "Trevor Cooper"},
		{Name: "Diana Richardson"},
		{Name: "Blake Watson"},
		{Name: "Courtney Brooks"},
		{Name: "Garrett Kelly"},
		{Name: "Alexis Howard"},
		{Name: "Colin Ward"},
		{Name: "Jasmine Torres"},
		{Name: "Seth Peterson"},
		{Name: "Gabrielle Gray"},
		{Name: "Evan Ramirez"},
		{Name: "Chloe James"},
		{Name: "Lucas Watson"},
		{Name: "Paige Bennett"},
	}

	// Check if readers already exist to avoid duplicates
	var count int64
	db.Model(&entity.Reader{}).Count(&count)
	if count > 0 {
		log.Println("Readers already seeded, skipping...")
		return nil
	}

	// Create readers in batches
	batchSize := 20
	for i := 0; i < len(readers); i += batchSize {
		end := i + batchSize
		if end > len(readers) {
			end = len(readers)
		}

		if err := db.Create(readers[i:end]).Error; err != nil {
			log.Printf("Error seeding readers batch %d-%d: %v", i, end-1, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d readers", len(readers))
	return nil
}
