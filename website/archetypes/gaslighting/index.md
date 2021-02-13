---
title: {{ replace .Name "-" " " | title }}
date: {{ dateFormat "2006-01-02" .Date }}
archive: _
img: screenshot.jpg
corpo: _
_build:
  render: never
---

