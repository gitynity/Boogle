import React from 'react'

function Listings( { books } ) {
  return (
	<ul>{
		books.forEach(book => {
			return <li>{book.name}<br/>{book.description}</li>
		})
	}</ul>
  )
}

export default Listings