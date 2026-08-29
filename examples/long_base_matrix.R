# Long Base-R transpiler test

normalize_vector <- function(x) {
  valid <- x[!is.na(x)]

  if (length(valid) == 0) {
    return(x)
  }

  minimum <- min(valid)
  maximum <- max(valid)

  if (maximum == minimum) {
    return(rep(0, length(x)))
  }

  (x - minimum) / (maximum - minimum)
}

classify_value <- function(x, lower, upper) {
  if (is.na(x)) {
    return("missing")
  }

  if (x < lower) {
    return("low")
  } else if (x > upper) {
    return("high")
  }

  "medium"
}

moving_average <- function(x, window = 3) {
  result <- rep(NA_real_, length(x))

  if (window <= 0 || window > length(x)) {
    return(result)
  }

  for (i in window:length(x)) {
    indices <- (i - window + 1):i
    values <- x[indices]
    result[i] <- mean(values, na.rm = TRUE)
  }

  result
}

matrix_summary <- function(x) {
  row_totals <- rowSums(x)
  column_totals <- colSums(x)

  list(
    rows = nrow(x),
    columns = ncol(x),
    row_totals = row_totals,
    column_totals = column_totals,
    transpose = t(x),
    diagonal = diag(x)
  )
}

fibonacci <- function(n) {
  if (n <= 0) {
    return(numeric(0))
  }

  if (n == 1) {
    return(0)
  }

  result <- numeric(n)
  result[1] <- 0
  result[2] <- 1

  if (n >= 3) {
    for (i in 3:n) {
      result[i] <- result[i - 1] + result[i - 2]
    }
  }

  result
}

describe_numeric <- function(x) {
  clean <- x[!is.na(x)]

  list(
    length = length(x),
    missing = sum(is.na(x)),
    minimum = min(clean),
    maximum = max(clean),
    mean = mean(clean),
    median = median(clean),
    variance = var(clean),
    standard_deviation = sd(clean),
    quantiles = quantile(clean, c(0, 0.25, 0.5, 0.75, 1))
  )
}

values <- c(12, 7, NA, 18, 4, 25, 9, 16, NA, 11)
weights <- seq(0.5, 5, length.out = length(values))
identifiers <- paste0("sample_", seq_along(values))

normalized <- normalize_vector(values)
averages <- moving_average(values, 3)

categories <- sapply(
  values,
  classify_value,
  lower = 8,
  upper = 17
)

results <- data.frame(
  id = identifiers,
  value = values,
  weight = weights,
  normalized = normalized,
  moving_average = averages,
  category = categories,
  stringsAsFactors = FALSE
)

results$weighted_value <- results$value * results$weight
results$is_missing <- is.na(results$value)
results$is_even <- !is.na(results$value) & results$value %% 2 == 0

cat("Original data:\n")
print(results)

cat("\nRows without missing values:\n")
complete_results <- results[complete.cases(results), ]
print(complete_results)

cat("\nRows ordered by value:\n")
ordered_results <- results[order(results$value, na.last = TRUE), ]
print(ordered_results)

category_counts <- table(results$category)
cat("\nCategory counts:\n")
print(category_counts)

numeric_description <- describe_numeric(values)
cat("\nNumeric description:\n")
print(numeric_description)

test_matrix <- matrix(
  seq_len(16),
  nrow = 4,
  ncol = 4,
  byrow = TRUE
)

rownames(test_matrix) <- paste0("row_", seq_len(nrow(test_matrix)))
colnames(test_matrix) <- paste0("column_", seq_len(ncol(test_matrix)))

cat("\nMatrix:\n")
print(test_matrix)

summary_matrix <- matrix_summary(test_matrix)

cat("\nMatrix summary:\n")
print(summary_matrix)

cat("\nMatrix multiplication:\n")
matrix_product <- test_matrix %*% t(test_matrix)
print(matrix_product)

cat("\nFibonacci sequence:\n")
fib <- fibonacci(20)
print(fib)

fib_even <- fib[fib %% 2 == 0]
fib_unique <- unique(fib)
fib_reverse <- rev(fib)

cat("\nEven Fibonacci values:\n")
print(fib_even)

cat("\nUnique Fibonacci values:\n")
print(fib_unique)

cat("\nReversed Fibonacci values:\n")
print(fib_reverse)

words <- c(
  "  Apple",
  "banana ",
  "Cherry",
  "apple",
  "BANANA",
  "date",
  "elderberry"
)

clean_words <- trimws(words)
lower_words <- tolower(clean_words)
upper_words <- toupper(clean_words)

cat("\nClean strings:\n")
print(clean_words)

cat("\nLowercase strings:\n")
print(lower_words)

cat("\nUppercase strings:\n")
print(upper_words)

cat("\nSorted unique strings:\n")
print(sort(unique(lower_words)))

joined <- paste(clean_words, collapse = ", ")
cat("\nJoined strings:\n")
print(joined)

split_result <- strsplit(joined, ", ")
cat("\nSplit strings:\n")
print(split_result)

number_text <- paste0(
  identifiers,
  "=",
  ifelse(is.na(values), "NA", as.character(values))
)

cat("\nFormatted values:\n")
print(number_text)

nested_list <- list(
  metadata = list(
    title = "Base R transpiler test",
    version = 1,
    active = TRUE
  ),
  data = results,
  matrix = test_matrix,
  fibonacci = fib,
  strings = clean_words
)

cat("\nNested-list names:\n")
print(names(nested_list))

cat("\nMetadata:\n")
print(nested_list$metadata)

accumulator <- 0
counter <- 1

while (counter <= 10) {
  if (counter == 5) {
    counter <- counter + 1
    next
  }

  accumulator <- accumulator + counter

  if (accumulator > 35) {
    break
  }

  counter <- counter + 1
}

cat("\nWhile-loop result:\n")
print(accumulator)

repeat_counter <- 0

repeat {
  repeat_counter <- repeat_counter + 1

  if (repeat_counter >= 3) {
    break
  }
}

cat("\nRepeat-loop result:\n")
print(repeat_counter)

apply_result <- apply(test_matrix, 1, function(row) {
  sum(row ^ 2)
})

lapply_result <- lapply(
  split(results$value, results$category),
  function(group) {
    group <- group[!is.na(group)]

    if (length(group) == 0) {
      return(NA_real_)
    }

    mean(group)
  }
)

cat("\nApply result:\n")
print(apply_result)

cat("\nGrouped means:\n")
print(lapply_result)

stopifnot(length(fib) == 20)
stopifnot(nrow(results) == length(values))
stopifnot(ncol(test_matrix) == 4)
stopifnot(all(dim(matrix_product) == c(4, 4)))
stopifnot(accumulator > 0)

cat("\nAll Base-R checks completed successfully.\n")