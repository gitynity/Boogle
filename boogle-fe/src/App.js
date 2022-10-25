import { Switch, Route } from 'react-router-dom';
import Home from './components/Home';
import NotFound from './components/NotFound';

function App() {
  return (
    <div className="App">
		<Switch>
			<Route path="/" exact>
				<Home />
			</Route>
			<Route path="*">
           	 	<NotFound />
          	</Route>
			  
		</Switch>
	
    </div>
  );
}

export default App;
