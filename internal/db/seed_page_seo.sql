INSERT INTO page_seo (
    id,
    page_path,
    page_name,
    meta_title,
    meta_description,
    meta_keywords,
    og_title,
    og_description,
    robots,
    created_at,
    updated_at
)
VALUES
    (
        gen_random_uuid(),
        '/',
        'Home',
        'Learn Islam Online | Iqra Initi',
        'Study Quran, Islamic studies, and essential Muslim knowledge through online courses with Iqra Initi.',
        'Islamic courses, Quran, online Islamic learning',
        'Learn Islam Online | Iqra Initi',
        'Explore online Quran and Islamic studies courses for every stage of learning.',
        'index,follow',
        NOW(),
        NOW()
    ),
    (
        gen_random_uuid(),
        '/courses',
        'Courses',
        'Islamic Courses Online | Iqra Initi',
        'Browse online Quran and Islamic studies courses taught for clear, practical learning.',
        'Islamic courses, Quran courses, Islamic studies',
        'Islamic Courses Online | Iqra Initi',
        'Find an online course for your Quran and Islamic studies journey.',
        'index,follow',
        NOW(),
        NOW()
    ),
    (
        gen_random_uuid(),
        '/about',
        'About',
        'About Iqra Initi | Learn and Grow',
        'Learn about Iqra Initi and our mission to make trusted Islamic learning accessible online.',
        'about Iqra Initi, Islamic learning, online education',
        'About Iqra Initi',
        'Our mission is to make trusted Islamic learning accessible to more people.',
        'index,follow',
        NOW(),
        NOW()
    ),
    (
        gen_random_uuid(),
        '/contact',
        'Contact',
        'Contact Iqra Initi',
        'Get in touch with Iqra Initi for questions about courses, learning, or your account.',
        'contact Iqra Initi, course support',
        'Contact Iqra Initi',
        'Our team is ready to help with course and account questions.',
        'index,follow',
        NOW(),
        NOW()
    ),
    (
        gen_random_uuid(),
        '/blog',
        'Blog',
        'Islamic Learning Articles | Iqra Initi',
        'Read articles and reflections on Quran, Islamic knowledge, and learning with Iqra Initi.',
        'Islamic articles, Quran, Islamic knowledge',
        'Islamic Learning Articles | Iqra Initi',
        'Explore articles on Quran, Islamic knowledge, and lifelong learning.',
        'index,follow',
        NOW(),
        NOW()
    )
ON CONFLICT (page_path) DO NOTHING;