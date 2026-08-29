square <- function(x) {
  x * x
}

values <- c(1, 2, 3, 4)
print(square(values))

total <- 0
for (value in values) {
  total <- total + value
}

total
