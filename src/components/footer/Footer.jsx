import React, { Component } from 'react';
import { connect } from 'react-redux';

import './Footer.scss';

const lightGrey = '#EDEDED';

export class Footer extends Component {
	
	constructor(props) {
		super(props);
	}

	render() {
		return (
			<div className='footer'>
				<div className='footer_content container'>
					<h2>Contact</h2>
					<p>Drop me a line</p>
					<hr />
					<div className='footer_content_social'>
						<ul>
							<li><a href='https://de.linkedin.com/in/felix-ortmann' target='_blank'><span className='icon'></span>https://de.linkedin.com/in/felix-ortmann</a></li>
							<li><a href='https://www.xing.com/profile/Felix_Ortmann2' target='_blank'><span className='icon'></span>https://www.xing.com/profile/Felix_Ortmann2</a></li>
							<li><a href='https://github.com/0ortmann' target='_blank'><span className='icon'></span>https://github.com/0ortmann</a></li>
							<li><a href='https://facebook.com/fixelexpress' target='_blank'><span className='icon'></span>https://facebook.com/fixelexpress</a></li>
							<li><a href='mailto:0ortmann@informatik.uni-hamburg.de' target='_blank'><span className='icon'></span>0ortmann@informatik.uni-hamburg.de</a></li>
						</ul>
					</div>
					<hr />
					<div className='footer_content_links'>
						<ul>
							<li><a href='/imprint'>{this.props.footer.imprintLink}</a></li>
							<li><a href='/privacy-policy'>{this.props.footer.privPolicyLink}</a></li>
						</ul>
					</div>
				</div>
			</div>
		);
	}
}

function mapStateToProps(state) {
	return {
		footer: state.lang.footer,
		contact: state.lang.address
	}
}

export default connect(mapStateToProps)(Footer);