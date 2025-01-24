package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"social/internal/store"

	"math/rand"
)

// Store some Dummy Data for the runnig query and optimize the query and use indexing.
var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe", "andrew", "bella", "carlton", "diana", "elijah", "faye", "gabriel", "harmony",
	"ian", "julia", "kenneth", "lorelei", "matthew", "nadia", "olivia", "paul",
	"quentin", "renee", "samuel", "theresa", "ursula", "valerie", "william", "xena",
	"yasmine", "zane", "albert", "brittany", "chris", "derek", "emily", "felix",
	"gavin", "heather", "isaac", "jasmine", "keira", "luke", "molly", "naomi",
	"orlando", "piper", "quincy", "roger", "sophia", "travis", "ulysses", "veronica",
	"wayne", "xander", "yvette", "zoey",
}

var titles = []string{
	"The Power of Habit", "Embracing Minimalism", "Healthy Eating Tips",
	"Travel on a Budget", "Mindfulness Meditation", "Boost Your Productivity",
	"Home Office Setup", "Digital Detox", "Gardening Basics",
	"DIY Home Projects", "Yoga for Beginners", "Sustainable Living",
	"Mastering Time Management", "Exploring Nature", "Simple Cooking Recipes",
	"Fitness at Home", "Personal Finance Tips", "Creative Writing",
	"Mental Health Awareness", "Learning New Skills", "Mastering Public Speaking", "Secrets of Effective Communication", "Living a Balanced Life",
	"Budget-Friendly Interior Design", "Understanding Emotional Intelligence", "Hiking Adventures for Beginners",
	"How to Start a Blog", "The Art of Storytelling", "Modern Home Automation Ideas",
	"Crafting the Perfect Resume", "Urban Gardening Tips", "Energy Conservation at Home",
	"Overcoming Procrastination", "The Beauty of Night Sky Photography", "Comfort Food Recipes",
	"Daily Meditation Practices", "Investment Basics for Beginners", "Exploring Local History",
	"Improving Sleep Quality", "Finding Hobbies That Inspire You",
}

var contents = []string{
	"In this post, we'll explore how to develop good habits that stick and transform your life.",
	"Discover the benefits of a minimalist lifestyle and how to declutter your home and mind.",
	"Learn practical tips for eating healthy on a budget without sacrificing flavor.",
	"Traveling doesn't have to be expensive. Here are some tips for seeing the world on a budget.",
	"Mindfulness meditation can reduce stress and improve your mental well-being. Here's how to get started.",
	"Increase your productivity with these simple and effective strategies.",
	"Set up the perfect home office to boost your work-from-home efficiency and comfort.",
	"A digital detox can help you reconnect with the real world and improve your mental health.",
	"Start your gardening journey with these basic tips for beginners.",
	"Transform your home with these fun and easy DIY projects.",
	"Yoga is a great way to stay fit and flexible. Here are some beginner-friendly poses to try.",
	"Sustainable living is good for you and the planet. Learn how to make eco-friendly choices.",
	"Master time management with these tips and get more done in less time.",
	"Nature has so much to offer. Discover the benefits of spending time outdoors.",
	"Whip up delicious meals with these simple and quick cooking recipes.",
	"Stay fit without leaving home with these effective at-home workout routines.",
	"Take control of your finances with these practical personal finance tips.",
	"Unleash your creativity with these inspiring writing prompts and exercises.",
	"Mental health is just as important as physical health. Learn how to take care of your mind.",
	"Learning new skills can be fun and rewarding. Here are some ideas to get you started.",
	"Public speaking can be intimidating, but with practice, anyone can become a confident speaker.",
	"Good communication is key to healthy relationships. Learn how to express yourself clearly.",
	"Discover the importance of balancing work, family, and personal time for a fulfilling life.",
	"Interior design doesn't have to be expensive. Here's how to refresh your home on a budget.",
	"Emotional intelligence is a crucial skill for both personal and professional growth.",
	"Hiking is a great way to stay fit and explore nature. Learn how to prepare for your first hike.",
	"Blogging can be a creative outlet or a career. Start your blog with these simple steps.",
	"Storytelling is a powerful tool for connection. Here are some techniques to captivate your audience.",
	"Home automation can simplify your life. Discover the latest gadgets and trends.",
	"Stand out to employers with a polished and professional resume using these tips.",
	"Transform your urban space into a lush garden with these beginner-friendly ideas.",
	"Energy conservation not only saves money but also helps the environment. Learn how to start.",
	"Procrastination is a common challenge. Find out how to overcome it and stay productive.",
	"Capture the beauty of the night sky with these photography tips for beginners.",
	"Try these comforting recipes that are perfect for relaxing weekends.",
	"Daily meditation can bring peace and clarity to your life. Here's how to make it a habit.",
	"Investing doesn't have to be complicated. Start with these beginner-friendly tips.",
	"Every town has a story. Learn how to uncover and appreciate the history of your local area.",
	"Sleep is vital for health and productivity. Find out how to improve your sleep quality.",
	"Discover hobbies that bring joy and creativity into your life with these practical ideas.",
}

var tags = []string{
	"Self Improvement", "Minimalism", "Health", "Travel", "Mindfulness",
	"Productivity", "Home Office", "Digital Detox", "Gardening", "DIY",
	"Yoga", "Sustainability", "Time Management", "Nature", "Cooking",
	"Fitness", "Personal Finance", "Writing", "Mental Health", "Learning",
	"Communication", "Life Balance", "Design", "Emotions", "Adventure",
	"Blogging", "Storytelling", "Technology", "Career", "Gardening",
	"Energy", "Productivity", "Photography", "Cooking", "Meditation",
	"Finance", "History", "Sleep", "Hobbies", "Inspiration",
}

var comments = []string{
	"Great post! Thanks for sharing.",
	"I completely agree with your thoughts.",
	"Thanks for the tips, very helpful.",
	"Interesting perspective, I hadn't considered that.",
	"Thanks for sharing your experience.",
	"Well written, I enjoyed reading this.",
	"This is very insightful, thanks for posting.",
	"Great advice, I'll definitely try that.",
	"I love this, very inspirational.",
	"Thanks for the information, very useful.",
	"This really spoke to me, thank you for sharing!",
	"I've been struggling with this; your advice is so helpful.",
	"Such a thoughtful post, I learned a lot!",
	"I'll definitely be implementing these tips.",
	"This is exactly what I needed to read today.",
	"Your perspective is refreshing, thank you for this.",
	"These ideas are brilliant, keep them coming!",
	"I tried this, and it worked wonderfully—thank you!",
	"Such a relatable post, you're not alone in this!",
	"Looking forward to more content like this!",
}

func Seed(store store.Storage, db *sql.DB) {

	ctx := context.Background()

	users := generateUsers(200)

	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, user, tx); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(400, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(1000, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			UserName: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			// Password: "123123",
			Password: store.Password{
				Hash: []byte("123123"),
				Text: "123123",
			},
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
