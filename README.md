# Vid serve
This is a toy implementation of a file service designed to securely return video files to a client.

It is purely designed to demonstrate the approach requireing minimal extra processing when scrubbing through the video 
timeline and has not been tested in any capacity, nor is it production quality

## Request Flow
```mermaid
sequenceDiagram
    participant Client
    participant ApplicationServer
    participant FileService

    Client->>ApplicationServer: Request page
    ApplicationServer->>FileService: Request access token for file
    FileService-->>ApplicationServer: Return token
    ApplicationServer-->>Client: Return page with tokens attached to file links
    Client->>FileService: Request file/chunk
    FileService-->>Client: Return file data
```
