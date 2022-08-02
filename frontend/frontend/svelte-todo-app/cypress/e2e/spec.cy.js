describe('home page', () => {
  it('h2 contains right text', () => {
    cy.visit("127.0.0.1:5000")
    cy.get("h2").contains("Svelte Todo App")
  })
  it('last todo contains "new todo" text', () => {
    cy.visit("127.0.0.1:5000")
    cy.get(".todo-input").type("new todo").type('{enter}')
    cy.get('.todo-item-label:last').contains('new todo')
    cy.get('.remove-item:last').click()
  })
  it('delete added todo', () => {
    cy.get(".todo-input").clear().type("first todo").type('{enter}')
    cy.get(".todo-input").clear().type("delete todo").type('{enter}')
    cy.get('.todo-item-label:last').should('have.text', 'delete todo')
    cy.get('.remove-item:last').click()
    cy.get('.todo-item-label:last').should('have.text', 'first todo')
    cy.get('.remove-item:last').click()

  })
  it('update todo status', () => {
    cy.get(".todo-input").clear().type("update todo").type('{enter}')
    cy.get('.todo-item-label:last').should('have.text', 'update todo')
    cy.get('[type="checkbox"]').last().click()
    cy.get('[type="checkbox"]').last().should('be.checked')
    cy.get('[type="checkbox"]').last().click()
    cy.get('[type="checkbox"]').last().should('not.be.checked')
    cy.get('.remove-item:last').click()
  })
})