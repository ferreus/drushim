import { JobsPage } from './app.po';

describe('jobs App', () => {
  let page: JobsPage;

  beforeEach(() => {
    page = new JobsPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!!');
  });
});
