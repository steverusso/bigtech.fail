---
title: {{ replace .Name "-" " " | title }}
date: {{ dateFormat "2006-01-02" .Date }}
image: /img/[the_event_image]
corpos: [ {{ index (split .Name "-") 0 }} ]
tags: [  ]
sources:
 - [ '', '' ]
---

