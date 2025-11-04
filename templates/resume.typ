#import "@local/modern-cv:0.9.0": *

#show link: set text(rgb(0, 102, 204))

#show: resume.with(
  author: (
    firstname: "Fangsong",
    lastname: "Long",
    email: "longfangsong981013@gmail.com",
    homepage: "https://longfangsong.github.io/en/",
    phone: "(+46) 76-451 43 68",
    github: "longfangsong",
    twitter: "longfangsong",
    birth: "Oct. 13th, 1998",
    linkedin: "fangsong-long-2b9b081aa",
    positions: (
      "Software Engineer",
      "Embedded developer",
      "Fullstack developer",
      "Cloud developer",
      "System software developer",
    ),
  ),
  keywords: ("Software Engineer"),
  description: "Fangsong complete resume",
  language: "en",
  colored-headers: true,
  show-footer: false,
  show-address-icon: true,
  paper-size: "us-letter",
  profile-picture: none,
  date: datetime.today().display(),
)

= Summary

A software developer with extensive knowledge. Previously worked at #link("https://www.pingcap.com/")[PingCAP, Inc.], a leading database company in China, where I contributed to the development of the transaction subsystem of a DBMS using Rust and Golang. Additionally, I am an active open-source contributor, both creating my own open-source projects, ranging from real-time OS kernels to RISC-V microcontrollers, and participating in the development of large open-source projects such as rust-analyzer and Rust HALs.

= Education

#resume-entry(
  title: [#link("https://www.chalmers.se/")[Chalmers University of Technology]],
  location: "Göteborg, Sweden",
  date: "Aug. 2023 - Aug. 2025",
  description: "Master of Software engineering and technology",
)
#resume-item[
  Got full grade on courses including:
  - Principles of Concurrent Programming
  - Real time systems
  - Formal methods in software development 
  - Quality assurance and testing
  - Language based security
  - Computability
]

#resume-entry(
  title: [#link("https://apply.shu.edu.cn/")[Shanghai University]],
  location: "Shanghai, China",
  date: "Sept. 2017 ‑ Aug. 2021",
  description: "Bachelor of computer science and  engineering",
)

#resume-item[
  - Technical director of #link("https://shuosc.github.io/")[Shanghai University Open Source Community]. Lead the community to develop several projects ranged from internal information platform to outsourcing projects from other departments of the university.
  - Works for the University’s Information Technology Office. Rewrite and maintain the OAuth system of the university.
]


#resume-entry(
  title: "Self Study and Online Courses",
  description: [Taking notes on my #link("https://longfangsong.github.io/en/")[Blog]],
)

#resume-item[
  - Teaching myself compiler techniques, mathematical logic, automata theory, type theory, physics of semiconductors, and many other things.
]

= Working Experience

#resume-entry(
  title: "Software Engineer",
  location: "Shanghai, China",
  date: "Jun. 2020 ‑ Sept. 2021 (Intern) , Sept. 2021 ‑ Sept 2022",
  description: [#link("https://www.pingcap.com/")[PingCAP, Inc.]],
)

#resume-item[
  - Take part in developing #link("https://www.pingcap.com/")[TiDB] & #link("https://tikv.org/")[TiKV], one of the most successful distributed databases around the world, focus mainly on the transaction part, worked with Golang and Rust.
  - Design, develop, and test the feature #link("https://docs.pingcap.com/tidb/dev/troubleshoot-lock-conflicts")[Lock View]. Help the user and the support engineer troubleshoot lock issues in the database system.
  - Refactor TiKV’s transaction #link("https://github.com/tikv/tikv/blob/571b5a263c7e84c2ab8aeb5feaebc8d50cae48cb/src/storage/txn/commands/mod.rs#%23L72")[command system] with macro. Increase the readability and maintainability of the code.
  - Develop #link("https://longfangsong.github.io/tipedia/en/index.html")[Tipedia], an unofficial knowledge base for helping developers to look for unfamiliar concepts in the system.
]

#resume-entry(
  title: "Software Engineer",
  location: "Shanghai, China",
  date: "Jun. 2018 ‑ Sept. 2018",
  description: [#link("https://www.synyi.com/")[SYNYI AI, Inc.]],
)

#resume-item[
  - Working on developing infrastructure for AI assisted medical diagnosis.
  - Develop #link("https://github.com/synyi/poplar")[Poplar], a high performance, web based open source annotation library for natural language processing needs.\
    Increase the first paint performance by 10x.\ Now this library is used by not only the company itself but also some other leading AI medical health companies in China including Alihealth.
]


= Projects

#resume-entry(
  title: [Extending RISC‑V ISA to Optimize Convolution Operation in AIoT Scenery],
  date: "Jan. - May 2021",
  description: "Bachelor Graduation Project",
)

#resume-item[
  - Benchmark space and time consumption of a simple convolution operation on a RISC-V CPU using different instruction set architecture module combinations in the #link("https://www.gem5.org/")[gem5] simulator.
  - Design a custom ISA module for convolution operation.
  - Extend RISC‑V GNU Compiler Toolchain to generate binary code which uses the new commands.
  - Extend gem5 simulator to simulate and benchmark the result binary.
  - Binary size decreased by around 50%, performance increased by a factor of 10 when using the new commands to implement convolution.
]

#resume-entry(
  title: [SHU RISC‑V suite],
  date: "Sept. 2020 - Now",
  description: "A RISC‑V MCU and compiling toolchain designed for it",
)
#resume-item[
- #link("https://github.com/shuosc/shuorv")[shuorv]. A simple RISC‑V MCU implementation in chisel, supports RV32I and basic interrupt handling, basic support of GPIO and Serial. Can boot FreeRTOS.
- #link("https://github.com/shuosc/come")[Come]. A toy C like language and it’s compiler tool chain, include a frontend, an optimizer, a RISC-V backend, a wasm backend, a RISC-V assembler and linker.
- We can compile a simple program written in Come to RISC‑V asm and binary code, and run it on shuorv.
]

#resume-entry(
  title: [A toy smart watch],
  location: github-link("longfangsong/watch"),
  date: "2021 - 2023",
  // description: "A RISC‑V MCU and compiling toolchain designed for it",
)
#resume-item[
- Design and assemble a PCB'A integrating an MCU, Wi-Fi module, LCD screen, magnetometer, and a pulse oximeter & heart rate sensor.

- Develope firmware in Rust, implementing features such as pulse monitoring and a pedometer.

- Achieve 20 FPS 240x320 video playback on a low-frequency MCU (STM32F103, which max frequency is 72MHz).
]



#resume-entry(
  title: [#link("https://lassvenska.pages.dev/")[Läss]],
  location: [#github-link("longfangsong/lass")],
  date: "May 2024 - Present",
  description: "A swedish learning platform",
)

#resume-item[
  - A Swedish learning platform built around articles from #link("https://sverigesradio.se/radioswedenpalattsvenska")[Radio Sweden på lätt svenska]. It allows users to read articles, click on unfamiliar words, and save them to a personal wordbook for review.
  - Leverages generative AI to explain word meanings and automatically create vocabulary and grammar exercises.
  - A local-first Progressive Web App (PWA) developed with Next.js and deployed on #link("https://www.cloudflare.com/")[Cloudflare]. Utilizes Google Gemini as the LLM provider, taking full advantage of edge computing and serverless architecture for *zero* operating costs.
  - Well received and actively used by my SFI classmates.
]
#resume-entry(
  title: "SHUHelper",
  location: github-link("shuosc/shuhelper"),
  date: "Sep. 2019 - Dec. 2022",
  description: "An information service platform for university students",
)

#resume-item[
  - Supports displaying school calendar, remind the student when and which classroom they should go for next class, record assignments, etc.
  - Use micro services architect. Deployed on a cloud Kubernetes cluster.
  - Widely used among students.
]
#resume-entry(
  title: [#link("https://github.com/SHUReeducation/autoAPI/wiki")[autoAPI]],
  location: github-link("SHUReeducation/autoAPI"),
  date: "Jun. 2021 - Dec. 2022",
  description: "A low‑code CRUD API generating tool",
)

#resume-item[
  - Read configure from configure file, SQL migrating file and/or meta tables in MySQL or PostgreSQL
  - Generating a standardized Golang RESTful micro service, which is ready to build with Docker and deployed with Kubernetes.
]
#resume-entry(
  title: "And more ...",
  description: "Experimental projects"
)

#resume-item[
  - Real time operating system kernel for embedded systems in #link("https://github.com/longfangsong/stm32-os")[C] and #link("https://github.com/longfangsong/rs-rtt")[Rust]
  - #link("https://github.com/shuosc/HPermission")[A high performance authentication gateway]
  - #link("https://github.com/longfangsong/tcpjunk")[A TCP package sniffer]
  - Even more on my #link("https://github.com/longfangsong/")[GitHub]
]

= Skills

#resume-skill-item(
  "Languages",
  (strong("C/C++"), strong("Javascript/Typescript"), strong("Golang"), strong("Rust"), strong("Python"), "Assembly (RISC-V & ARMv7)", "Verilog" , "C#", "Erlang", "Scala (Chisel DSL)", "SQL", "HTML", "CSS", "WASM", "Kotlin", "Lua", "GDScript", "Erlang"),
)

#resume-skill-item(
  "Frontend",
  (strong("React"), strong("Next.js"), strong("TailwindCSS"), strong("shadcn/ui"), strong("Vue"), strong("Angular"), "Svelte"),
)

#resume-skill-item(
  "Backend",
  (strong("Django"), strong("fastapi"), strong("Gin"), "axum", "beego", "actix"),
)

#resume-skill-item("Databases", (strong("PostgreSQL (Neon/On-premises)"), strong("MySQL"), strong("SQLite"), strong("Redis"), "MongoDB"))

#resume-skill-item("Cloud Platforms", ("Cloudflare", "AWS", "Azure", "Google Cloud Platform", "Vercel", "Supabase"))

#resume-skill-item("Dev Ops", ("Kubernetes", "GitHub Actions", "ArgoCD", "GitLab CI/CD", "Docker", "Terraform", "Ansible"))

#resume-skill-item("AI/ML library", ("Pytorch", "OpenCV"))

#resume-skill-item("LLM related", ("ollama", "LangChain", "n8n"))

#resume-skill-item("Misc", ("Linux", "Shell" , "Git", "PCB Design", "3D modeling", "3D printing", "SEO", "Prometheus", "Grafana", "Agile", "TDD", "Domain-Driven Design"))

#resume-skill-item("Spoken Languages", (strong("Chinese"), strong("English"), strong("Svenska (B1+, SFI nivå D)"), "Einfaches Deutsch (A2)"))

= Interests

#resume-skill-item(
  "Technical Writing",
  ([I write blogs to record what I have learned.],),
)

#resume-skill-item("Open Source", ([I have contributed to #link("https://github.com/rust-lang/rust-analyzer")[rust‑analyzer], #link("https://github.com/cloudflare/next-on-pages")[next-on-pages], #link("https://chaos-mesh.org/")[chaos‑mesh], #link("https://github.com/riscv-rust/k210-hal")[k210‑hal], etc.], ))

#resume-skill-item("Electronic DIY", ([I have built projects like automatic light, solar powered power bank, e‑paper reader, etc.],))
