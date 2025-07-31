#!/usr/bin/env node
declare class ReadingTrackerMCPServer {
    private server;
    private openLibraryApi;
    private readingTrackerApi;
    constructor();
    private setupToolHandlers;
    private handleSearchBooksByTitle;
    private handleSearchAuthorsByName;
    private handleGetBookById;
    private handleGetAuthorInfo;
    private handleGetBookCover;
    private handleGetAuthorPhoto;
    private handleEnrichBookData;
    run(): Promise<void>;
}
export { ReadingTrackerMCPServer };
//# sourceMappingURL=index.d.ts.map