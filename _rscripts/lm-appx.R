#!/usr/bin/Rscript
# author: ggian
# date: Mon Oct 24 11:40:15 EEST 2016
#
# Script used to create a linear model used for regression


args <- commandArgs(trailingOnly=TRUE)
if (length(args) < 2 ) { 
	cat("Please provide the training set and test sets.\n") 
	quit()
}

datafile <- args[1]
testfile <- args[2]
train.data <- read.csv(datafile)
test.data <- read.csv(testfile)

fit <- lm(class ~ ., data=train.data)
pred <- predict(fit, test.data)
#cat(sum((pred-test.data$class)^2)/length(test.data))
for (i in 1:length(pred)) {
#		cat(test.data$class[i])
#		cat("\t")
		cat(pred[i])
		cat("\n")
}
