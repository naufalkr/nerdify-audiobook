export const mockBooks = [
  {
    book_id: "1",
    title: { String: "Pride and Prejudice" },
    author: { String: "Jane Austen" },
    image_url: { String: "https://archive.org/services/img/pride_prejudice_0711_librivox" },
    genre: "Romance"
  },
  {
    book_id: "2", 
    title: { String: "The Adventures of Sherlock Holmes" },
    author: { String: "Arthur Conan Doyle" },
    image_url: { String: "https://archive.org/services/img/adventures_sherlock_holmes_1112_librivox" },
    genre: "Mystery"
  },
  {
    book_id: "3",
    title: { String: "Alice's Adventures in Wonderland" },
    author: { String: "Lewis Carroll" },
    image_url: { String: "https://archive.org/services/img/alices_adventures_wonderland_1001_librivox" },
    genre: "Children's Fiction"
  },
  {
    book_id: "4",
    title: { String: "Frankenstein" },
    author: { String: "Mary Shelley" },
    image_url: { String: "https://archive.org/services/img/frankenstein_0902_librivox" },
    genre: "Horror & Supernatural Fiction"
  },
  {
    book_id: "5",
    title: { String: "The Time Machine" },
    author: { String: "H.G. Wells" },
    image_url: { String: "https://archive.org/services/img/time_machine_1010_librivox" },
    genre: "Science Fiction"
  }
]

export const mockChapters = {
  "1": [
    {
      Title: "Chapter 1: It is a truth universally acknowledged",
      Link: "https://archive.org/download/pride_prejudice_0711_librivox/prideandprejudice_01_austen_64kb.mp3"
    },
    {
      Title: "Chapter 2: Mr. Bennet was among the earliest",
      Link: "https://archive.org/download/pride_prejudice_0711_librivox/prideandprejudice_02_austen_64kb.mp3"
    },
    {
      Title: "Chapter 3: Not all that Mrs. Bennet",
      Link: "https://archive.org/download/pride_prejudice_0711_librivox/prideandprejudice_03_austen_64kb.mp3"
    }
  ],
  "2": [
    {
      Title: "A Scandal in Bohemia",
      Link: "https://archive.org/download/adventures_sherlock_holmes_1112_librivox/adventuressherlockholmes_01_doyle_64kb.mp3"
    },
    {
      Title: "The Red-Headed League",
      Link: "https://archive.org/download/adventures_sherlock_holmes_1112_librivox/adventuressherlockholmes_02_doyle_64kb.mp3"
    }
  ],
  "3": [
    {
      Title: "Chapter 1: Down the Rabbit Hole",
      Link: "https://archive.org/download/alices_adventures_wonderland_1001_librivox/alice_wonderland_01_carroll_64kb.mp3"
    },
    {
      Title: "Chapter 2: The Pool of Tears",
      Link: "https://archive.org/download/alices_adventures_wonderland_1001_librivox/alice_wonderland_02_carroll_64kb.mp3"
    }
  ],
  "4": [
    {
      Title: "Letter 1",
      Link: "https://archive.org/download/frankenstein_0902_librivox/frankenstein_01_shelley_64kb.mp3"
    },
    {
      Title: "Letter 2",
      Link: "https://archive.org/download/frankenstein_0902_librivox/frankenstein_02_shelley_64kb.mp3"
    }
  ],
  "5": [
    {
      Title: "Chapter 1: Introduction",
      Link: "https://archive.org/download/time_machine_1010_librivox/timemachine_01_wells_64kb.mp3"
    },
    {
      Title: "Chapter 2: The Machine",
      Link: "https://archive.org/download/time_machine_1010_librivox/timemachine_02_wells_64kb.mp3"
    }
  ]
}

export const mockUser = {
  email: "demo@example.com",
  name: "Demo User",
  id: "demo-user-123"
}