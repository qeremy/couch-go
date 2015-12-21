package main

import (
    // test_uuid "./uuid"
    // test_client "./client"
    // test_server "./server"
    // test_database "./database"
    test_document "./document"
    // test_document_attachment "./document_attachment"
)

func main() {
    /* client */
    // test_client.TestClientResponseStatus()
    // ...

    /* server */
    // test_server.TestPing()
    // ...

    /* database */
    // test_database.TestPing()
    // ...

    /* document */
    // test_document.TestPing()
    // test_document.TestIsExists()
    // test_document.TestIsNotModified()
    // test_document.TestFind()
    // test_document.TestFindRevisions()
    // test_document.TestFindRevisionsExtended()
    // test_document.TestFindAttachments()
    test_document.TestSave()

    /* document attachment */
    // test_document_attachment.TestPing()
    // test_document_attachment.TestFind()
    // test_document_attachment.TestReadFile()
    // test_document_attachment.TestToArray()
    // test_document_attachment.TestToJson()
    // test_document_attachment.TestSave()
    // test_document_attachment.TestRemove()

    /* uuid */
    // test_uuid.TestAll()
}
