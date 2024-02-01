package errors

import "errors"

var PostIdMissingError = errors.New("Post Id is Invalid, please enter a value greater than 0")
var TitleMissingError = errors.New("Title is missing")
var ContentMissingError = errors.New("Content is missing")
var AuthorMissingError = errors.New("Author is missing")
var PublicationDateMissingError = errors.New("Publication Date is missing")
var InvalidPublicationDateError = errors.New("Publication Date is invalid, should be in the format dd-mm-yyyy")
var TagsMissingError = errors.New("Tags are missing")
