import React, { Component } from 'react';

import './Main.scss';

export default class Main extends Component {
    
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <main className='about'>
                <h1>FELIX ORTMANN</h1>
                <div className='about__subtitle'>
                    &mdash; Student, Entwickler, Verrückter &mdash;
                </div>
                <img className='about__img'/>
                <div className='about__description'>
                    <p>
                        Baue gerne verteilte Systeme und mag Herausforderungen. Am liebsten bin ich auf Linux Servern unterwegs,
                        eine Shell ist alles was ich brauche. Handwerke immer wieder in allen möglichen Gebieten wie Web oder Bigdata 
                        und lerne begeistert neue Dinge.
                    </p>
                    <p>
                        Momentan studiere ich an der Uni Hamburg und bin im zweiten Master Semester. Ich Entwickle Software seit 2011, 
                        es begann mit einem Forschungsprojekt zur verteilten Petrinetz Simulation an der Uni. In der Wirtschaft habe ich 
                        mich früh in Richtung Web bewegt und gleite nun zunehmend in DevOps Bereiche über.
                    </p>
                    <p>
                        Ich bin gerne in der Natur unterwegs, entweder mit Rucksack &amp; Zelt oder ich mache mit der Angel jagd auf heimische Süßwasserräuber. 
                        Bei schlechtem Wetter bastel ich an kleinen Projekten, sowie etwa <a href='https://github.com/0ortmann/wg-tools'>wg-tools</a>, 
                        das gerade <a href='https://wg-tools.de'>online geganen</a> ist! Ein paar andere nette Kleinigkeiten und nützliche Dinge sind auf meinem Github-Account zu finden.
                    </p>
                </div>
            </main>
        );
    }
}